package main

import (
	"log"
	"time"
)

func main() {

	serverAddr := [4]byte{0x7f, 0x00, 0x00, 0x01} // 127.0.0.1
	serverPort := 2222

	server := TFOServer{ServerAddr: serverAddr, ServerPort: serverPort}
	err := server.Bind()
	if err != nil {
		log.Fatalln("Failed to bind socket:", err)
	}

	// Create a new routine ("thread") and wait for connection from client
	go server.Accept()

	client := TFOClient{ServerAddr: serverAddr, ServerPort: serverPort}

	err = client.Send()
	if err != nil {
		log.Fatalln("Failed to send to server:", err)
	}

	// Give the server a chance to output and exit
	time.Sleep(100 * time.Millisecond)

}
