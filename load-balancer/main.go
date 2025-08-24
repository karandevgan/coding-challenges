package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"sync/atomic"
)

var backendServers = []string{"localhost:8080", "localhost:8081", "localhost:8082"}
var backendServerIndex int32 = 0

func main() {
	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Printf("Error listening on TCP port 80: %s\n", err)
		return
	}
	defer listener.Close()
	log.Printf("Listening on TCP port 80\n")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s\n", err)
			continue
		}
		log.Printf("Received request from %s\n", conn.RemoteAddr())
		sIndex := atomic.LoadInt32(&backendServerIndex)
		log.Printf("Current backend server index: %d\n", sIndex)
		go handleConnection(conn, sIndex)
		atomic.CompareAndSwapInt32(&backendServerIndex, sIndex, (sIndex+1)%int32(len(backendServers)))
	}
}

func handleConnection(conn net.Conn, sIndex int32) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in handleConnection: %v", r)
			conn.Close()
		}
	}()
	reader := bufio.NewReader(conn)
	bServer := backendServers[sIndex]
	log.Printf("Forwarding request to backend server %s\n", bServer)
	dConn, err := net.Dial("tcp", bServer)
	if err != nil {
		log.Printf("Error connecting to backend server: %s\n", err)
		conn.Close()
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic in copy to backend: %v", r)
				conn.Close()
				dConn.Close()
			}
		}()
		_, err := io.Copy(dConn, reader)
		if err != nil {
			log.Printf("Error copying data to backend server: %s\n", err)
		}
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic in copy from backend: %v", r)
				conn.Close()
				dConn.Close()
			}
		}()
		_, err := io.Copy(conn, dConn)
		if err != nil {
			log.Printf("Error copying data from backend server: %s\n", err)
		}
	}()
}
