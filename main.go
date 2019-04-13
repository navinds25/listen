package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
	"syscall"

	"github.com/syossan27/tebata"
)

var (
	address  string
	protocol string
	ports    string
)

func listen(address, port, protocol string, t *tebata.Tebata) error {
	listenAddress := address + ":" + port
	fmt.Println(listenAddress, protocol)
	lis, err := net.Listen(protocol, listenAddress)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer lis.Close()
	t.Reserve(lis.Close)
	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("from: %v to: %s\n", conn.RemoteAddr(), port)
	}
}

func main() {
	flag.StringVar(&address, "address", "0.0.0.0", "Address to listen on")
	flag.StringVar(&ports, "ports", "8121,8122,8123", "Ports to listen on, comma separated, no spaces.")
	flag.StringVar(&protocol, "protocol", "tcp", "Protocol use")
	flag.Parse()
	t := tebata.New(syscall.SIGINT, syscall.SIGTERM)
	splitPorts := strings.Split(ports, ",")
	lastPort := splitPorts[len(splitPorts)-1]
	splitPorts = splitPorts[:len(splitPorts)-1]
	for _, port := range splitPorts {
		go func(port string) {
			if err := listen(address, port, protocol, t); err != nil {
				fmt.Println("Error: ", err)
			}
		}(port)
		//go listen(address, port, protocol, t)
	}
	if err := listen(address, lastPort, protocol, t); err != nil {
		fmt.Println(err)
	}
}
