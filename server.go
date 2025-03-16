package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"tcpserver/commands"
)

type Server struct {
	port        int
	loggedUsers map[string]bool
}

func NewServer(port int) *Server {
	return &Server{
		port:        port,
		loggedUsers: map[string]bool{},
	}
}

func (s *Server) Start() {

	// Start listening on the specified port
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	fmt.Println("Server ready for incoming connections...")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Accept incoming connections (in async, so that we can accept many)
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {

	fmt.Println("Connection established, waiting for commands...")
	for {

		cmd, cmdErr := commands.ParseCommand(conn)
		if cmdErr == io.EOF {
			fmt.Println("Client disconnected")
			break
		}
		if cmdErr != nil {
			fmt.Println("Error while parsing command:", cmdErr)
			break
		}

		resp, procErr := cmd.Process()
		if procErr != nil {
			fmt.Println("Error while processing command:", procErr)
			break
		}

		resp.Write(conn)
	}
}
