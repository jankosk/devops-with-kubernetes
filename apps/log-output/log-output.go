package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

func main() {
	id := uuid.NewString()
	port := ":8080"
	var currentTimestamp atomic.Pointer[string]

	go startLoggingTicker(&currentTimestamp, id, time.Second*5)

	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		timestampPtr := currentTimestamp.Load()
		if timestampPtr != nil {
			fmt.Fprintf(writer, "%s: %s\n", *timestampPtr, id)
		}
	})

	fmt.Printf("Server listening on port %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}

func startLoggingTicker(currentTimestamp *atomic.Pointer[string], id string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for t := range ticker.C {
		newTimestamp := t.Format(time.RFC3339)
		currentTimestamp.Store(&newTimestamp)
		fmt.Printf("%s: %s\n", newTimestamp, id)
	}
}
