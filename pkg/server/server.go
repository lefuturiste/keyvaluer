package server

import (
	"fmt"
	"net"
	"os"
)

// StartAsync - Start the TCP server with a go routine
func StartAsync(port string, pass string) {
	if port != "0" && port != "" {
		os.Setenv("PORT", port)
	}
	os.Setenv("REQUIRED_PASS", pass)
	go serverLoop()
}

// Start - Start the TCP server
func Start() {
	serverLoop()
}

func serverLoop() {
	var host string = os.Getenv("HOST")
	var port string = os.Getenv("PORT")
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}
	l, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + host + ":" + port)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleClient(conn)
	}
}
