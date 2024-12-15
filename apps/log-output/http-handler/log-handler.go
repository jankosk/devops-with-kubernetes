package main

import (
	"bufio"
	"dwk/common"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var logsPath = common.GetEnv("LOGS_PATH", "/tmp")
var configPath = common.GetEnv("CONFIG_PATH", "/tmp")
var message = common.GetEnv("MESSAGE", "")
var pingPongUrl = common.GetEnv("PING_PONG_URL", "http://localhost:3001")

func main() {
	port := ":3000"

	http.HandleFunc("/", handleLogRequest)

	log.Printf("Server listening on port %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}

func handleLogRequest(w http.ResponseWriter, r *http.Request) {
	logLine, err := readLastLine(filepath.Join(logsPath, "logs.txt"))
	if err != nil {
		common.HandleErr(w, "Unable to read log file", http.StatusInternalServerError, err)
		return
	}

	pingPongCount, err := fetchPingPongs()
	if err != nil {
		common.HandleErr(w, "Unable to fetch ping pongs", http.StatusInternalServerError, err)
		return
	}

	configText, err := readConfigFile()
	if err != nil {
		common.HandleErr(w, "Unable to read config file", http.StatusInternalServerError, err)
		return
	}

	fmt.Fprintf(w, "file content: %s\nenv variable: %s\n%s\nPing / Pongs: %d\n", configText, message, logLine, pingPongCount)
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

func readConfigFile() (string, error) {
	data, err := os.ReadFile(filepath.Join(configPath, "information.txt"))
	return string(data), err
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
