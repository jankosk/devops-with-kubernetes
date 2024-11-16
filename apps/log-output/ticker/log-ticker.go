package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func blockForever() {
	select {}
}

func main() {
	id := uuid.NewString()
	logFilePath := os.Getenv("LOGS_PATH")
	f, err := os.Create(filepath.Join(logFilePath, "logs.txt"))
	check(err)

	go startLoggingTicker(f, id, time.Second*5)

	blockForever()
}

func startLoggingTicker(f *os.File, id string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("Ticker with %vs interval started\n", interval.Seconds())
	for t := range ticker.C {
		newTimestamp := t.Format(time.RFC3339)
		f.WriteString(fmt.Sprintf("%s: %s\n", newTimestamp, id))
	}
}
