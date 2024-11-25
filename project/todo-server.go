package main

import (
	"dwk/common"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var PORT = common.GetEnv("PORT", "8080")
var filesPath = common.GetEnv("FILES_PATH", "/tmp")

func main() {
	port := ":" + PORT

	http.Handle("GET /", http.FileServer(http.Dir("public")))

	http.HandleFunc("GET /random-image", randomImageHandler)

	fmt.Printf("Server listening on port %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
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
