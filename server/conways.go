package main

import (
	"github.com/tarm/goserial"
	"io"
	"log"
	"os"
	"math/rand"
	"time"
)

const xboardsize int = 28
const yboardsize int = 24

type Board struct {
	cur []bool // t/f alive/dead
	xboards int
	yboards int
	boardsize int
}

func NewBoard(xboards int, yboards int) Board {
	boardsize := (yboards*yboardsize) * (xboards*xboardsize)
	b := Board{
		make([]bool, boardsize),
		xboards,
		yboards,
		boardsize,
	}
	return b
}

func (b Board) ToCoords(i int) (x,y int) {
	y = i / (xboardsize * b.xboards)
	x = (i % (xboardsize * b.xboards))

	return
}

func (b Board) ToInd(x, y int) int {
	return (y * (xboardsize * b.xboards)) + x
}

type Change struct {
	x,y int
	dir bool // t/f on/off
}

func Bytes(c Change) (b []byte) {
	xb := c.x / xboardsize // x board number
	// TODO: Implement y board integration
	// yb := c.y / yboardsize // y board number
	b = make([]byte, 3)
	b[0] = byte(c.x % xboardsize) // x coord
	b[1] = byte(c.y % yboardsize) // y coord
	b[2] = (byte(xb) << 1)
	if c.dir {
		b[2] |= 0x1
	}

	return
}

type Game struct {
	Board
	Changeset []Change
	ToChange int
	serial io.ReadWriteCloser
}

func NewGame(xboards int, yboards int, serial io.ReadWriteCloser) Game {
	b := NewBoard(xboards, yboards)
	c := make([]Change, b.boardsize)
	return Game{
		NewBoard(xboards, yboards),
		c,
		0,
		serial,
	}
}

func (g Game) Update() {
	for i := 0; i < g.ToChange; i++ {
		c := g.Changeset[i]
		g.cur[g.ToInd(c.x, c.y)] = c.dir
		g.serial.Write(Bytes(c))
		//os.Stdout.Sync()
		//fmt.Fprintln(os.Stderr, g.ToInd(c.x, c.y), "\tx,", c.x," y,",
		//c.y, ":\n\t", Bytes(c))
	}
	g.ToChange = 0
}

func (g Game) Scramble() {
	for i := 0; i < g.boardsize; i++ {
		r := rand.Float32()
		if r > 0.5 {
			c := &(g.Changeset[g.ToChange])
			c.x, c.y = g.ToCoords(i)
			c.dir = true
			g.ToChange++
		}
	}
	g.Update()
}

// Flipdot Functions
func Clear(state bool) []byte {
	if(state) {
		return []byte{0x00, 0x00, 0xF0}
	} else {
		return []byte{0x00, 0x00, 0xE0}
	}
}

func main() {
	log.Println("Life begins")

	// Open serial port
	log.Print("Opening Serial Port... ")
	c := &serial.Config{Name: os.Args[1], Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done")

	log.Print("Sending Clear... ")
	// Clear
	_, err = s.Write([]byte{0x00, 0x00, 0xF0})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Sent")

	// Seed randomness
	rand.Seed(time.Now().UnixNano())

	// Create game
	g := NewGame(2, 1, s)

	time.Sleep(1 * time.Second)
	// Scramble
	g.Scramble()

	for {}
}
