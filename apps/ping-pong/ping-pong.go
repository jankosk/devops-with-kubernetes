package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

func main() {
	var counter atomic.Uint32
	port := ":8080"

	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(writer, "pong %d\n", counter.Load())
		counter.Add(1)
	})

	fmt.Printf("Server listening on port %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
