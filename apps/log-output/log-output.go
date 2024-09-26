package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func main() {
	uuid := uuid.NewString()
	for {
		timestamp := time.Now().Format(time.RFC3339)
		logMessage := fmt.Sprintf("%s: %s", timestamp, uuid)
		fmt.Println(logMessage)
		time.Sleep(5 * time.Second)
	}
}
