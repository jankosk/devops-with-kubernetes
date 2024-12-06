package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

func main() {
	var counter atomic.Uint32
	port := ":3001"

	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		counter.Add(1)

		fmt.Fprintf(writer, "%d\n", counter.Load())
	})

	fmt.Printf("Server listening on port %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
