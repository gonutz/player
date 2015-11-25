package main

import (
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	http.HandleFunc("/", serve)
	http.ListenAndServe(":8080", nil)
}

func serve(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func listDir() []file {
	var files []file
	filepath.Walk(workingDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			files = append(files, file{path, info.IsDir()})
		}
		if info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	})
	return files
}

type file struct {
	path  string
	isDir bool
}

func activatePath(f file) {
	if f.isDir {
		workingDir = f.path
	} else {
		// TODO play movie
	}
}
