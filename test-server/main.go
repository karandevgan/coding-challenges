package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := "8080" // pass port as argument
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	server := &http.Server{
		Addr: ":" + port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Connection", "keep-alive")
			fmt.Fprintf(w, "Hello from server on port %s\n", port)
		}),
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(server.ListenAndServe())
}
