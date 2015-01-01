package main

import (
	"log"
	"time"
)

func main() {
	serverFD, err := serverBind()
	if err != nil {
		log.Fatalln("Failed to bind socket:", err)
	}

	go serverAccept(serverFD)

	err = clientSend()

	if err != nil {
		log.Fatalln("Failed to send to server:", err)
	}

	// Wait for response and exit
	time.Sleep(1 * time.Second)

}
