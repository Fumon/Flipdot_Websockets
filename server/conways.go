package main

import (
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/tarm/goserial"
)

// Default set to my large boards
const xboardsize int = 28
const yboardsize int = 24

// Conways rules
//
// If Alive
const underpop = 2 // Less than = dead
const staymin = 2
const staymax = 3
const overpop = 3 // Greater than = dead
// If Dead
const repro = 3 // Equal to = alive

type Board struct {
	cur       []bool // t/f alive/dead
	xboards   int
	yboards   int
	boardsize int
}

func NewBoard(xboards int, yboards int) Board {
	boardsize := (yboards * yboardsize) * (xboards * xboardsize)
	b := Board{
		make([]bool, boardsize),
		xboards,
		yboards,
		boardsize,
	}
	return b
}

func (b Board) ToCoords(i int) (x, y int) {
	y = i / (xboardsize * b.xboards)
	x = (i % (xboardsize * b.xboards))

	return
}

func (b Board) ToInd(x, y int) int {
	return (y * (xboardsize * b.xboards)) + x
}

// Get the value of a coordinate, returning false when outside bounds.
func (b Board) GetCoord(x, y int) bool {
	if x < 0 || x >= b.xboards*xboardsize || y < 0 || y >
		b.yboards*yboardsize {
		return false
	}
	return b.GetInd(b.ToInd(x, y))
}

func (b Board) GetInd(index int) bool {
	if index < 0 || index >= b.boardsize {
		return false
	}

	return b.cur[index]
}

func (b Board) SetInd(index int, value bool) {
	b.cur[index] = value
}

// Get the surrounding 8 cells of an index
func (b Board) GetSurrounding(ind int) []bool {
	// Convert to coordinates
	x, y := b.ToCoords(ind)

	return []bool{
		b.GetCoord(x-1, y-1), // Above
		b.GetCoord(x, y-1),
		b.GetCoord(x+1, y-1),
		b.GetCoord(x-1, y+1), // Below
		b.GetCoord(x, y+1),
		b.GetCoord(x+1, y+1),
		b.GetCoord(x-1, y), // Left
		b.GetCoord(x+1, y), // Right
	}
}

type Change struct {
	index int
	dir   bool // t/f on/off
}

func Bytes(x, y int, dir bool) (b []byte) {
	xb := x / xboardsize // x board number
	// TODO: Implement y board integration
	b = make([]byte, 3)
	b[0] = byte(x % xboardsize) // x coord
	b[1] = byte(y % yboardsize) // y coord
	b[2] = (byte(xb) << 1)
	if dir {
		b[2] |= 0x1
	}

	return
}

type Game struct {
	Board
	Changeset []Change
	ToChange  int
	serial    io.ReadWriteCloser
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

func (g *Game) AddChange(index int, dir bool) {
	c := &(g.Changeset[g.ToChange])
	c.index = index
	c.dir = dir
	g.ToChange++
}

func (g *Game) Update() {
	log.Println(g.ToChange, " Changes")
	for i := 0; i < g.ToChange; i++ {
		c := g.Changeset[i]
		g.SetInd(c.index, c.dir)
		// Convert to coordinates
		x, y := g.ToCoords(c.index)
		bytes := Bytes(x, y, c.dir)
		g.serial.Write(bytes)
	}
	g.ToChange = 0

	g.serial.Write([]byte{0x00, 0x00, 0xD0}) // Ack request
	// Wait for ack
	buf := make([]byte, 1)
	for {
		_, err := g.serial.Read(buf)
		if err != nil {
			if err == io.EOF {
				continue
			} else {
				log.Fatalln("Error reading from device: ", err)
			}
		}

		if buf[0] == 'H' {
			break
		}
	}

}

func (g *Game) Scramble() {
	for i := 0; i < g.boardsize; i++ {
		r := rand.Float32()
		if r > 0.5 {
			g.AddChange(i, true)
		}
	}
}

// Takes a boolean array of 8 cell's and the target cell's alive status and
// returns whether the target cell should live or die.
func (g *Game) returnAlive(target bool, set []bool) bool {
	// Sum the alive cells
	sum := 0
	for _, cell := range set {
		if cell {
			sum++
		}
	}

	// Compare with rules
	if target {
		if sum < underpop {
			return false
		} else if sum >= staymin && sum <= staymax {
			return true
		} else if sum > overpop {
			return false
		}
	} else {
		if sum == repro {
			return true
		} else {
			return false
		}
	}

	return false
}

// Stages one step of game of life in the changeset
func (g *Game) Step() {
	for i := 0; i < g.boardsize; i++ {
		// Get current state
		cur := g.GetInd(i)

		n := g.returnAlive(cur, g.GetSurrounding(i))

		if n != cur { // Create a change
			g.AddChange(i, n)
		}
	}
}

// Flipdot Functions
func Clear(state bool) []byte {
	if state {
		return []byte{0x00, 0x00, 0xF0}
	} else {
		return []byte{0x00, 0x00, 0xE0}
	}
}

func main() {
	log.Println("Life begins")

	// Open serial port
	log.Print("Opening Serial Port... ")
	c := &serial.Config{Name: os.Args[1], Baud: 57600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done")

	log.Print("Sending Clear... ")
	// Clear
	_, err = s.Write(Clear(false))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Sent")

	// Seed randomness
	rand.Seed(time.Now().UnixNano())

	// Create game
	gp := NewGame(3, 1, s)
	g := &gp

	// Scramble
	g.Scramble()
	time.Sleep(4 * time.Second)

	log.Println("Sending Scramble")
	g.Update()

	gens := 0
	for {
		g.Step()
		g.Update()
		gens++
		log.Println("Gen -- ", gens)
	}

	for {
	}
}
