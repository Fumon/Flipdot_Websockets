package flipdot

const fontw = 6
const fonth = 7

type FlipdotText struct {
	*Flipdot
	Xpanels int
}

func (f *FlipdotText) ResetCursor() error {
	return f.writePacket(ResetCursorPkt[:])
}

func (f *FlipdotText) PlaceCursor(x, y int) error {
	xp := x / xboardsize
	nx := x % xboardsize
	return f.writePacket([]byte{byte(nx), byte(y), byte(xp), MoveCursor})
}

func (f *FlipdotText) Letter(a rune) error {
	return f.writePacket([]byte{byte(a), 0x00, 0x00, PrintCharacter})
}

func (f *FlipdotText) Print(b string) (int, error) {
	for i, x := range b {
		err := f.Letter(x)
		if err != nil {
			return i, err
		}
	}

	return len(b), nil
}

func (f *FlipdotText) PrintCentered(b string) error {
	slen := len(b) * fontw

	xoff := ((xboardsize * f.Xpanels) - slen) / int(2)
	err := f.PlaceCursor(xoff, (yboardsize-fonth)/int(2))
	if err != nil {
		return err
	}
	_, err = f.Print(b)
	return err
}
