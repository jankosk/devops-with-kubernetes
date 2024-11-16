package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var logFilePath = os.Getenv("LOGS_PATH")

func main() {
	port := ":8080"

	http.HandleFunc("/", handleLogRequest)

	fmt.Printf("Server listening on port %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}

func handleLogRequest(w http.ResponseWriter, r *http.Request) {
	logLine, err := readLastLine(filepath.Join(logFilePath, "logs.txt"))
	if err != nil {
		http.Error(w, "Unable to read log file", http.StatusInternalServerError)
		log.Printf("Error reading log file: %v\n", err)
		return
	}

	fmt.Fprintln(w, logLine)
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
