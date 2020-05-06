package server

import (
	"net"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// StartAsync - Start the TCP server with a go routine
func StartAsync(port string, pass string) {
	if port != "0" && port != "" {
		os.Setenv("PORT", port)
	}
	os.Setenv("REQUIRED_PASS", pass)
	os.Setenv("LOG_LEVEL", "debug")
	go serverLoop()
}

func initLogging() {
	var strLogLevel = os.Getenv("LOG_LEVEL")
	var logLevel log.Level
	if strLogLevel == "" || strLogLevel == "info" {
		logLevel = log.InfoLevel
	}
	switch strings.ToLower(strLogLevel) {
	case "trace":
		logLevel = log.TraceLevel
	case "debug":
		logLevel = log.DebugLevel
	case "warn":
		logLevel = log.WarnLevel
	case "error":
		logLevel = log.ErrorLevel
	case "fatal":
		logLevel = log.FatalLevel
	case "panic":
		logLevel = log.PanicLevel
	}
	log.SetLevel(logLevel)
}

// Start - Start the TCP server
func Start() {
	serverLoop()
}

func serverLoop() {
	initLogging()
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
		log.Error("Error listening: ", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	log.Info("Listening on " + host + ":" + port)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error("Error accepting connexion: ", err.Error())
			os.Exit(1)
		}
		go handleClient(conn)
	}
}
