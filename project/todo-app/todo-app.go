package main

import (
	"bytes"
	"dwk/common"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var PORT = common.GetEnv("PORT", "8080")
var filesPath = common.GetEnv("FILES_PATH", "/tmp")
var todoApiUrl = common.GetEnv("TODO_API_URL", "http://localhost:8083")

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func main() {
	port := ":" + PORT

	http.HandleFunc("GET /", indexPageHandler)
	http.HandleFunc("POST /", formPostHandler)
	http.HandleFunc("GET /random-image", randomImageHandler)

	log.Printf("Server listening on port %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("public/index.html")
	if err != nil {
		common.HandleErr(w, "Failed parsing index.html", http.StatusInternalServerError, err)
		return
	}
	todos, err := fetchTodos()
	if err != nil {
		common.HandleErr(w, "Failed fetching todos", http.StatusInternalServerError, err)
		return
	}

	data := TodoPageData{
		PageTitle: "My TODO list",
		Todos:     todos,
	}
	templ.Execute(w, data)
}

func formPostHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		common.HandleErr(w, "Failed to parse form data", http.StatusBadRequest, err)
		return
	}

	title := r.PostForm.Get("title")
	todo := Todo{Title: title, Done: false}
	if err := createTodo(todo); err != nil {
		common.HandleErr(w, "Failed to create todo", http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func randomImageHandler(w http.ResponseWriter, r *http.Request) {
	data, err := getCachedRandomImage()
	if err != nil {
		common.HandleErr(w, "Failed to get image", http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(data)
}

func getCachedRandomImage() ([]byte, error) {
	imgPath := filepath.Join(filesPath, "random.jpeg")
	fileInfo, err := os.Stat(imgPath)
	if err != nil || time.Since(fileInfo.ModTime()) > time.Hour {
		_, err := storeRandomImage(imgPath)
		if err != nil {
			return nil, err
		}
	}
	return os.ReadFile(imgPath)
}

func storeRandomImage(path string) (int64, error) {
	resp, err := http.Get("https://picsum.photos/1200")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	file, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return io.Copy(file, resp.Body)
}

func createTodo(todo Todo) error {
	todoJson, err := json.Marshal(todo)
	if err != nil {
		return fmt.Errorf("failed to serialize todo %w", err)
	}

	res, err := http.Post(todoApiUrl, "application/json", bytes.NewReader(todoJson))
	if err != nil {
		return fmt.Errorf("failed to send POST request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}

func fetchTodos() ([]Todo, error) {
	res, err := http.Get(todoApiUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch todos: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var todos []Todo
	if err := json.NewDecoder(res.Body).Decode(&todos); err != nil {
		return nil, fmt.Errorf("failed to decode todos JSON: %w", err)
	}

	return todos, nil
}
