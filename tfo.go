package main

import (
	"log"
	"net"
	"time"

	"github.com/droundy/goopt"
)

var connect = goopt.String([]string{"-s", "--server"}, "127.0.0.1", "Server to connect to (and listen if listening)")
var port = goopt.Int([]string{"-p", "--port"}, 2222, "Port to connect to (and listen to if listening)")

var listen = goopt.Flag([]string{"-l", "--listen"}, []string{}, "Create a listening TFO socket", "")

func main() {

	goopt.Parse(nil)

	var serverAddr [4]byte

	IP := net.ParseIP(*connect)
	copy(serverAddr[:], IP[12:16])

	if *listen {

		server := TFOServer{ServerAddr: serverAddr, ServerPort: *port}
		err := server.Bind()
		if err != nil {
			log.Fatalln("Failed to bind socket:", err)
		}

		// Create a new routine ("thread") and wait for connection from client
		go server.Accept()

	}

	client := TFOClient{ServerAddr: serverAddr, ServerPort: *port}

	err := client.Send()
	if err != nil {
		log.Fatalln("Failed to send to server:", err)
	}

	// Give the server a chance to output and exit
	time.Sleep(100 * time.Millisecond)

}
