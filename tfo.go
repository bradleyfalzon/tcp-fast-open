package main

import (
	"bytes"
	"log"
	"net"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/droundy/goopt"
)

var connect = goopt.String([]string{"-s", "--server"}, "127.0.0.1", "Server to connect to (and listen if listening)")
var port = goopt.Int([]string{"-p", "--port"}, 2222, "Port to connect to (and listen to if listening)")

var listen = goopt.Flag([]string{"-l", "--listen"}, []string{}, "Create a listening TFO socket", "")

func main() {

	goopt.Parse(nil)

	// IPv4 only for no real reason, could be v6 by adjusting the sizes
	// here and where it's used
	var serverAddr [4]byte

	IP := net.ParseIP(*connect)
	if IP == nil {
		log.Fatal("Unable to process IP: ", *connect)
	}

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

	// Give the server a chance to receive, process the packet and print results
	time.Sleep(100 * time.Millisecond)

	success, cached, err := checkTcpMetrics(*connect)
	if err != nil {
		log.Println("ip tcp_metrics failure:", err)
	} else {
		var response string
		if success {
			response = "TFO success to IP " + *connect
		} else {
			response = "TFO failure to IP " + *connect
		}
		if len(cached) > 0 {
			response += " " + strings.Join(cached, ", ")
		}
		log.Println(response)
	}

}

// Use `ip tcp_metrics` to check whether we received a cookie or not. Only
// available in later versions of iproute
func checkTcpMetrics(ip string) (success bool, cached []string, err error) {

	cmd := exec.Command("ip", "tcp_metrics", "show", ip)

	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return
	}

	reFOc := regexp.MustCompile(" fo_cookie ([a-z0-9]+)")
	reFOmss := regexp.MustCompile(" fo_mss ([0-9]+)")
	reFOdrop := regexp.MustCompile(" fo_syn_drops ([0-9./]sec ago)")

	cookie := reFOc.FindStringSubmatch(out.String())
	mss := reFOmss.FindStringSubmatch(out.String())
	drop := reFOdrop.FindStringSubmatch(out.String())

	success = len(cookie) > 0

	if len(cookie) > 0 {
		cached = append(cached, "cookie: "+cookie[1])
	}

	if len(mss) > 0 {
		cached = append(cached, "mss: "+mss[1])
	}

	if len(drop) > 0 {
		cached = append(cached, "syn_drops: "+drop[1])
	}

	return

}
