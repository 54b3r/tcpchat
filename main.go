package main

import (
	"log"
	"net"

	srv "github.com/54b3r/tcpchat/pkg"
)

func errLog(err error) {
	if err != nil {
		log.Fatalf("unable  to start server: %s", err.Error())
	}
}

func main() {
	s := srv.NewServer()
	// start our server in a go routine
	go s.Run()

	listener, err := net.Listen(s.Protocol, ":"+s.Port)
	if err != nil {
		srv.Logger(true, "[ERROR]: unable to start server: %s", err.Error())
	}

	defer listener.Close()
	srv.Logger(false, "[INFO]: Started Server on port %s", s.Port)

	// continuous for loop to accept incomming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			// If there is an error, we will continue to allow new connections and handle failed connections as well
			srv.Logger(false, "[ERROR]: Unable to accept connection %s", err.Error())
			continue
		}
		go s.NewClient(conn)
	}
}
