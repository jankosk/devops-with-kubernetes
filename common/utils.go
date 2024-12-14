package common

import (
	"log"
	"net/http"
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

func HandleErr(w http.ResponseWriter, message string, status int, err error) {
	if err != nil {
		log.Printf("Error: %s - %v", message, err)
	} else {
		log.Printf("Error: %s", message)
	}
	http.Error(w, message, status)
}
