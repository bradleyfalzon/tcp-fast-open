package main

import (
	"errors"
	"log"
	"syscall"
)

const MSG_FASTOPEN int = 0x20000000

func clientSend() (err error) {

	log.Println("Creating client socket...")

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return
	}
	defer syscall.Close(fd)

	log.Println("Got client socket FD:", fd)

	sa := &syscall.SockaddrInet4{Port: 2222, Addr: [4]byte{127, 0, 0, 1}}

	//data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	data := []byte("HELLO TCP FO")

	log.Printf("Sending to server: %#v\n", data)

	err = syscall.Sendto(fd, data, MSG_FASTOPEN, sa)
	if err != nil {
		if err == syscall.EOPNOTSUPP {
			return errors.New("TCP Fast Open client support is unavailable (unsupported kernel or disabled, see /proc/sys/net/ipv4/tcp_fastopen).")
		}
		log.Println("Got error sending to server:", err)
		return
	}

	return
}
