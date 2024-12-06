package main

import (
	"bufio"
	"dwk/common"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var logsPath = common.GetEnv("LOGS_PATH", "/tmp")
var pingPongUrl = common.GetEnv("PING_PONG_URL", "http://localhost:3001")

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

	pingPongCount, err := fetchPingPongs()
	if err != nil {
		http.Error(w, "Unable to fetch ping pongs", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s\nPing / Pongs: %d\n", logLine, pingPongCount)
}

func fetchPingPongs() (int, error) {
	res, err := http.Get(pingPongUrl)
	if err != nil {
		return 0, fmt.Errorf("error making GET request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %w", err)
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(body)))
	if err != nil {
		return 0, fmt.Errorf("error parsing response to integer: %w", err)
	}

	return count, nil
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
