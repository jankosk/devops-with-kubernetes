package main

import (
	"dwk/common"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func blockForever() {
	select {}
}

func main() {
	id := uuid.NewString()
	logFilePath := common.GetEnv("LOGS_PATH", "/tmp")
	f, err := os.Create(filepath.Join(logFilePath, "logs.txt"))
	common.CheckErr(err, "Failed to create log file")

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
