package main

import (
	"log"
	"os"
	"time"

	"./flipdot"
)

const xboards = 3

func main() {
	log.Println("Life begins")

	// Open serial port
	log.Print("Opening Flipdot on serialport ", os.Args[1], " ...")

	f, err := flipdot.NewFlipdotSerial(os.Args[1], 57600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.Println("Done")

	// Clear
	err = f.Clear(false)
	if err != nil {
		log.Fatal(err)
	}

	ft := &flipdot.FlipdotText{f, xboards}

	// Time loop
	after := time.After(time.Millisecond)
	for {
		t := <-after
		after = time.After(t.Truncate(time.Minute).Add(time.Minute).Sub(t))
		ft.PrintCentered(t.Format("Jan _2 15:04"))
	}
}
