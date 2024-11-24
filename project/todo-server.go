package main

import (
	"dwk/common"
	"fmt"
	"net/http"
)

var PORT = common.GetEnv("PORT", "8080")

func main() {
	port := ":" + PORT
	fmt.Printf("Server listening on port %s\n", port)

	http.Handle("/", http.FileServer(http.Dir("public")))

	http.ListenAndServe(port, nil)
}
