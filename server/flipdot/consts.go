package flipdot

// For the larger 24x28 boards
const (
	xboardsize = 28
	yboardsize = 24
)

// Commandbytes
// Leave the two lowest bits alone as it is used to indicate which direction and on which panel
// to flip a single dot in FlipOne mode.
const (
	FlipOne        = 0x80
	SetFontDim     = 0x90
	resetCursor    = 0xA0
	MoveCursor     = 0xB0
	PrintCharacter = 0xC0
	ackRequest     = 0xD0
	clearOff       = 0xE0
	clearOn        = 0xF0
)

// Some non-variable packets
var (
	ResetCursorPkt = [...]byte{0x00, 0x00, 0x00, resetCursor}
	AckRequestPkt  = [...]byte{0x00, 0x00, 0x00, ackRequest}
	ClearOffPkt    = [...]byte{0x00, 0x00, 0x00, clearOff}
	ClearOnPkt     = [...]byte{0x00, 0x00, 0x00, clearOn}
)
