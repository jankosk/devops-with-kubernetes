package common

import (
	"log"
	"os"
)

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func CheckErr(e error, msg string) {
	if e != nil {
		log.Fatalf("%s: %v", msg, e)
	}
}
