package main

import (
	"log"
	"os"

	"./flipdot"
)

// Default set to my large boards
const xboardsize int = 28
const yboardsize int = 24

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

	log.Print("Sending Clear... ")
	// Clear
	err = f.Clear(false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Sent")

	log.Println("Placing cursor")
	ft := &flipdot.FlipdotText{f, 3}

	log.Println("Printing ", os.Args[2])
	ft.PrintCentered(os.Args[2])
}
