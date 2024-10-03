package main

import (
	"fmt"
	"net/http"
	"os"
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

var PORT = getEnv("PORT", "8080")

func main() {
	port := ":" + PORT
	fmt.Printf("Server listening on port %s\n", port)

	http.Handle("/", http.FileServer(http.Dir("public")))

	http.ListenAndServe(port, nil)
}
