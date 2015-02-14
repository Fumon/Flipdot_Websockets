package flipdot

import (
	"io"

	"github.com/tarm/goserial"
)

type Flipdot struct {
	io.ReadWriteCloser
}

// NewFlipdotSerial returns a new *Flipdot using a serial port connection
func NewFlipdotSerial(port string, baud int) (*Flipdot, error) {
	c := &serial.Config{Name: port, Baud: baud}
	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}

	return &Flipdot{s}, nil
}

// Internal conveninence
func (f *Flipdot) writePacket(b []byte) error {
	_, err := f.Write(b)
	return err
}

// General Utilities

// Clear the connected board to on (true) or off (false)
func (f *Flipdot) Clear(state bool) error {
	err := error(nil)
	if state {
		err = f.writePacket(ClearOnPkt[:])
	} else {
		err = f.writePacket(ClearOffPkt[:])
	}
	return err
}

// SendAckRequest sends a request for an ack byte to be returned from the board
// once its work has been completed and it is ready for more commands.
func (f *Flipdot) SendAckRequest() error {
	return f.writePacket(AckRequestPkt[:])
}

// Flip flips one dot in x,y coordinates
func (f *Flipdot) Flip(x, y int, on bool) error {
	xp := x / xboardsize         // x panel number
	xout := byte(x % xboardsize) // x coord
	yout := byte(y % yboardsize) // y coord
	cbyte := (byte(xp) << 1)
	if on {
		cbyte |= 0x1
	}

	return f.writePacket([]byte{xout, yout, cbyte, FlipOne})
}
