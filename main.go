package main

import (
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	http.HandleFunc("/", serve)
	http.ListenAndServe(":8080", nil)
}

func serve(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(htmlList(listDir())))
}

func htmlList(list fileList) string {
	items := ""
	for _, f := range list {
		items += "\n<li>" + f.path + "</li>"
	}
	return `<html>
	<body>
		<ul style="list-style-type:none">` + items +
		`		</ul>
	</body>
</html>`
}

func listDir() fileList {
	var files fileList
	filepath.Walk(workingDir, func(path string, info os.FileInfo, err error) error {
		if path == workingDir {
			return nil
		}

		if err == nil {
			files = append(files, file{path, info.IsDir()})
		}
		if info != nil && info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	})
	sort.Sort(files)
	return files
}

type file struct {
	path  string
	isDir bool
}

type fileList []file

func (f fileList) Len() int { return len(f) }

func (f fileList) Less(i, j int) bool {
	if f[i].isDir != f[j].isDir {
		if f[i].isDir {
			return true
		}
		return false
	}

	return strings.Compare(strings.ToLower(f[i].path), strings.ToLower(f[j].path)) < 0
}

func (f fileList) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func activatePath(f file) {
	if f.isDir {
		workingDir = f.path
	} else {
		// TODO play movie
	}
}
