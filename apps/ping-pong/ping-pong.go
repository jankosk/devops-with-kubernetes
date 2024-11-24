package main

import (
	"dwk/common"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync/atomic"
)

func main() {
	var counter atomic.Uint32
	logPath := common.GetEnv("LOGS_PATH", "/tmp")
	port := ":3001"

	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		f, err := os.Create(filepath.Join(logPath, "ping-pongs.txt"))
		if err != nil {
			http.Error(writer, "Unable to read file", http.StatusInternalServerError)
			return
		}
		counter.Add(1)
		_, err = f.WriteString(fmt.Sprintf("%d", counter.Load()))
		if err != nil {
			http.Error(writer, "Unable to write file", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(writer, "%d\n", counter.Load())
	})

	fmt.Printf("Server listening on port %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
