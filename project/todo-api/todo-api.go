package main

import (
	"dwk/common"
	"encoding/json"
	"fmt"
	"net/http"
)

type Todo struct {
	Title string
	Done  bool
}

var todos []Todo = []Todo{}

func main() {
	port := ":" + common.GetEnv("PORT", "8083")

	http.HandleFunc("GET /", getTodosHandler)
	http.HandleFunc("POST /", createTodoHandler)

	fmt.Printf("Server listening on port %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}

func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, todos)
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	todos = append(todos, todo)

	respondWithJSON(w, http.StatusCreated, todo)
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
	}
}
