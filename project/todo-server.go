package main

import (
	"dwk/common"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var PORT = common.GetEnv("PORT", "8080")
var filesPath = common.GetEnv("FILES_PATH", "/tmp")

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

	http.HandleFunc("GET /random-image", randomImageHandler)

	fmt.Printf("Server listening on port %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("public/index.html")
	if err != nil {
		http.Error(w, "Failed parsing index.html", http.StatusInternalServerError)
		return
	}
	data := TodoPageData{
		PageTitle: "My TODO list",
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}
	templ.Execute(w, data)
}

func randomImageHandler(w http.ResponseWriter, r *http.Request) {
	data, err := getCachedRandomImage()
	if err != nil {
		http.Error(w, "Failed to get image", http.StatusInternalServerError)
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
