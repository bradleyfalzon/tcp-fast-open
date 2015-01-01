package main

import (
	"errors"
	"log"
	"syscall"
)

const TCP_FASTOPEN int = 23
const LISTEN_BACKLOG int = 23

func serverBind() (fd int, err error) {

	log.Println("Creating server socket...")

	fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		if err == syscall.ENOPROTOOPT {
			return fd, errors.New("TCP Fast Open server support is unavailable (unsupported kernel).")
		}
		return
	}

	log.Println("Got server socket FD:", fd)

	log.Println("Binding socket...")

	sa := &syscall.SockaddrInet4{Port: 2222, Addr: [4]byte{127, 0, 0, 1}}

	err = syscall.Bind(fd, sa)
	if err != nil {
		return
	}

	log.Println("Bound socket.")

	log.Println("Attempting to listening on socket...")

	err = syscall.Listen(fd, LISTEN_BACKLOG)
	if err != nil {
		return
	}

	log.Println("Listening on socket with fd: ", fd)

	log.Println("Setting socket options...")

	err = syscall.SetsockoptInt(fd, syscall.SOL_TCP, TCP_FASTOPEN, 1)
	if err != nil {
		return
	}

	log.Println("Set socket options.")

	return

}

func serverAccept(fd int) {

	log.Println("Waiting for connection on fd: ", fd)

	defer syscall.Close(fd)

	for {
		cFd, cSockaddr, err := syscall.Accept(fd)
		if err != nil {
			log.Fatalln("Failed to accept(): ", err)
		}

		go handleClient(cFd, cSockaddr.(*syscall.SockaddrInet4))

	}

}

func handleClient(fd int, cSockaddr *syscall.SockaddrInet4) {

	defer closeClient(fd)

	log.Printf("Connection received, new FD: %d, remote port: %d, remote ip: %d.%d.%d.%d\n",
		fd, cSockaddr.Port, cSockaddr.Addr[0], cSockaddr.Addr[1], cSockaddr.Addr[2], cSockaddr.Addr[3])

	// Handle request

	buf := make([]byte, 24)

	n, err := syscall.Read(fd, buf)
	if err != nil {
		log.Println("Failed to read() client:", err)
		return
	}

	log.Printf("Read %d bytes: %#v", n, string(buf[:n]))

}

func closeClient(fd int) {

	err := syscall.Shutdown(fd, syscall.SHUT_RDWR)
	if err != nil {
		log.Println("Failed to shutdown() connection:", err)
	}

	err = syscall.Close(fd)
	if err != nil {
		log.Println("Failed to close() connection:", err)
	}

}
