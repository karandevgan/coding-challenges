package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var backendServers = []string{"localhost:8080", "localhost:8081", "localhost:8082"}
var healthCheckUrl = "/"
var healthCheckTimeout = 5 * time.Second
var healthCheckInterval = 10 * time.Second
var connectionTimeout = 60 * time.Second
var responseTimeout = 60 * time.Second

var serverList atomic.Value // holds []string
var nextServerCounter uint64

func main() {
	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Printf("Error listening on TCP port 80: %s\n", err)
		return
	}
	defer listener.Close()
	log.Printf("Listening on TCP port 80\n")

	// Initialize snapshot
	serverList.Store(append([]string(nil), backendServers...))

	ticker := time.NewTicker(healthCheckInterval)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			log.Printf("Performing health check on backend servers\n")
			currentHealthyServers := make([]string, 0, len(backendServers))
			for _, server := range backendServers {
				if !checkServerHealth(server) {
					log.Printf("Removing server %s from backend servers\n", server)
				} else {
					log.Printf("Server %s is healthy\n", server)
					currentHealthyServers = append(currentHealthyServers, server)
				}
			}
			// Update snapshot
			serverList.Store(currentHealthyServers)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s\n", err)
			continue
		}
		currentHealthyServers := serverList.Load().([]string)
		if len(currentHealthyServers) == 0 {
			log.Printf("No healthy servers, closing connection\n")
			// Send 503 Service Unavailable response
			_, _ = conn.Write([]byte("HTTP/1.1 503 Service Unavailable\r\nContent-Length: 19\r\n\r\nService Unavailable\n"))
			_ = conn.Close()
			continue
		}
		idx := atomic.AddUint64(&nextServerCounter, 1) - 1
		nextServer := currentHealthyServers[idx%uint64(len(currentHealthyServers))]
		log.Printf("Received request from %s\n", conn.RemoteAddr())
		go handleConnection(conn, nextServer)
	}
}

func checkServerHealth(server string) bool {
	client := http.Client{
		Timeout: healthCheckTimeout,
	}
	res, err := client.Get("http://" + server + healthCheckUrl)
	if err != nil {
		log.Printf("Error performing health check on server %s: %s\n", server, err)
		return false
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		log.Printf("Health check on server %s failed with status code %d\n", server, res.StatusCode)
		return false
	}
	return true
}

func handleConnection(conn net.Conn, nextServer string) {
	defer conn.Close()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in handleConnection: %v", r)
		}
	}()
	_ = conn.SetDeadline(time.Now().Add(responseTimeout))
	reader := bufio.NewReader(conn)
	bServer := nextServer
	log.Printf("Forwarding request to backend server %s\n", bServer)

	// Dial backend server
	dConn, err := net.DialTimeout("tcp", bServer, connectionTimeout)
	if err != nil {
		log.Printf("Error connecting to backend server: %s\n", err)
		_, _ = conn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\nContent-Length: 15\r\n\r\nBad Gateway\n"))
		return
	}
	defer dConn.Close()
	_ = dConn.SetDeadline(time.Now().Add(responseTimeout))

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic in copy to backend: %v", r)
			}
			// Close write side of backend connection when client stops sending
			if tcpConn, ok := dConn.(*net.TCPConn); ok {
				_ = tcpConn.CloseWrite()
			}
			wg.Done()
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
			}
			// Close write side of client connection when backend stops sending
			if tcpConn, ok := conn.(*net.TCPConn); ok {
				_ = tcpConn.CloseWrite()
			}
			wg.Done()
		}()
		_, err := io.Copy(conn, dConn)
		if err != nil {
			log.Printf("Error copying data from backend server: %s\n", err)
		}
	}()

	// Wait for both directions to finish to avoid leaks
	wg.Wait()
}
