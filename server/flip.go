package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"net/http"
	"os"
)

// Websocket main
func FlipServ(ws *websocket.Conn) {
	//buff := make([3]byte)
	for {
		buff := make([]byte, 60)
		n, err := ws.Read(buff[:])
		if n != 3 || err != nil {
			log.Println("Problem with websocket content")
			break
		}
		log.Println("Recieved Content:\n\t ", buff[0:3])
		os.Stdout.Write(buff[0:3])
	}
}

func main() {
	log.Println("Life begins")
	http.Handle("/", http.FileServer(http.Dir("../gui/initializr")))
	http.Handle("/flipdot", websocket.Handler(FlipServ))
	log.Println("Starting server")
	err := http.ListenAndServe("0.0.0.0:7779", nil)
	if err != nil {
		log.Panicln("ListenAndServe error: ", err.Error())
	}
}
