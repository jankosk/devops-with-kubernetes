package main

import (
	"database/sql"
	"dwk/common"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Todo struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var (
	port        = common.GetEnv("PORT", "8083")
	db_host     = common.GetEnv("DB_HOST", "localhost")
	db_user     = common.GetEnv("DB_USERNAME", "app")
	db_password = common.GetEnv("DB_PASSWORD", "example")
)

func main() {
	var dbUrl = fmt.Sprintf("host=%s user=%s password=%s dbname=postgres sslmode=disable", db_host, db_user, db_password)
	db, err := sql.Open("postgres", dbUrl)
	common.CheckErr(err, "Failed to connect to the database")
	defer db.Close()

	err = initDb(db)
	common.CheckErr(err, "Failed to initialize database")

	http.HandleFunc("GET /", getTodosHandler(db))
	http.HandleFunc("POST /", createTodoHandler(db))

	log.Printf("Server listening on port %s\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}

func getTodosHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		todos, err := getTodos(db)
		if err != nil {
			common.HandleErr(w, "Failed to fetch todos", http.StatusBadRequest, err)
			return
		}
		respondWithJSON(w, http.StatusOK, todos)
	}
}

func createTodoHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo Todo
		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			common.HandleErr(w, "Invalid JSON payload", http.StatusBadRequest, err)
			return
		}
		todo, err := createTodo(db, todo)
		if err != nil {
			common.HandleErr(w, "Failed to create todo", http.StatusBadRequest, err)
			return
		}
		respondWithJSON(w, http.StatusCreated, todo)
	}
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		common.HandleErr(w, "Failed to serialize response", http.StatusInternalServerError, err)
	}
}

func initDb(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		done BOOLEAN NOT NULL DEFAULT FALSE
	)`)
	return err
}

func createTodo(db *sql.DB, parsedTodo Todo) (Todo, error) {
	var todo Todo
	query := `
	INSERT INTO todos (title, done)
	VALUES ($1, $2)
	RETURNING id, title, done`
	row := db.QueryRow(query, parsedTodo.Title, parsedTodo.Done)
	if err := row.Scan(&todo.Id, &todo.Title, &todo.Done); err != nil {
		return todo, err
	}
	return todo, nil
}

func getTodos(db *sql.DB) ([]Todo, error) {
	var todos []Todo = []Todo{}
	rows, err := db.Query(`SELECT * FROM todos`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		if err = rows.Scan(&todo.Id, &todo.Title, &todo.Done); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}
