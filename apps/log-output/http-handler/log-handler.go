package main

import (
	"bufio"
	"dwk/common"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var logsPath = common.GetEnv("LOGS_PATH", "/tmp")

func main() {
	port := ":3000"

	http.HandleFunc("/", handleLogRequest)

	fmt.Printf("Server listening on port %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}

func handleLogRequest(w http.ResponseWriter, r *http.Request) {
	logLine, err := readLastLine(filepath.Join(logsPath, "logs.txt"))
	if err != nil {
		http.Error(w, "Unable to read log file", http.StatusInternalServerError)
		return
	}
	pingPongs, err := readLastLine(filepath.Join(logsPath, "ping-pongs.txt"))
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}
	pingPongCount, err := strconv.Atoi(pingPongs)
	if err != nil {
		http.Error(w, "Unable to parse ping pongs count", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s\nPing / Pongs: %d\n", logLine, pingPongCount)
}

func readLastLine(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lastLine string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	return lastLine, nil
}
