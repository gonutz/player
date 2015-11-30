package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

var (
	player  videoPlayer = &omxplayer{}
	mux     sync.Mutex
	wdFiles fileList
)

func main() {
	player = &stubVideoPlayer{}
	listDir()
	http.HandleFunc("/", serve)
	http.ListenAndServe(":80", nil)
}

func serve(w http.ResponseWriter, r *http.Request) {
	if player.isRunning() {
		redirect := false
		if r.FormValue("pause") != "" {
			log(player.playPause())
			redirect = true
		}
		if r.FormValue("stop") != "" {
			log(player.stopVideo())
			redirect = true
		}
		if r.FormValue("lower") != "" {
			log(player.volumeDown())
			redirect = true
		}
		if r.FormValue("louder") != "" {
			log(player.volumeUp())
			redirect = true
		}
		if r.FormValue("backSmall") != "" {
			log(player.back30Seconds())
			redirect = true
		}
		if r.FormValue("forwardSmall") != "" {
			log(player.forward30Seconds())
			redirect = true
		}
		if r.FormValue("backBig") != "" {
			log(player.back10Minutes())
			redirect = true
		}
		if r.FormValue("forwardBig") != "" {
			log(player.forward10Minutes())
			redirect = true
		}
		if redirect {
			http.Redirect(w, r, "", http.StatusMovedPermanently)
		}
	} else {
		path, err := url.QueryUnescape(r.URL.Path)
		if len(path) < 3 || err != nil {
			log(err)
			http.Redirect(w, r, "/|"+url.QueryEscape(workingDir),
				http.StatusMovedPermanently)
		} else {
			path = path[2:]
			info, err := os.Stat(path)
			log(err)
			if err != nil {
				http.Redirect(w, r, "/|"+url.QueryEscape(workingDir),
					http.StatusMovedPermanently)
			} else {
				if info.IsDir() {
					workingDir = path
					listDir()
				} else {
					log(player.playVideo(path))
				}
			}
		}
	}

	if player.isRunning() {
		w.Write([]byte(videoControls))
	} else {
		w.Write([]byte(htmlList(wdFiles)))
	}
}

func htmlList(list fileList) string {
	items := ""
	for _, f := range list {
		suffix := ""
		if f.isDir {
			suffix = " ..."
		}
		items += "\n<li><a href=\"" + url.QueryEscape("|"+f.path) + "\">" + f.path +
			suffix + "</a></li>"
	}
	return `<html>
	<head>
		<style>
			body {
				font-size: 200%;
			}
		</style>
	</head>
	<body>
		<ul style="list-style-type:none">` + items +
		`		</ul>
	</body>
</html>`
}

func listDir() {
	wdFiles = nil
	filepath.Walk(workingDir, func(path string, info os.FileInfo, err error) error {
		if path == workingDir {
			return nil
		}

		if err == nil {
			wdFiles = append(wdFiles, file{path, info.IsDir()})
		}
		if info != nil && info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	})
	wdFiles = append(wdFiles, file{filepath.Dir(workingDir), true})
	sort.Sort(wdFiles)
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

const videoControls = `<html>
	<head>
		<style>
			input[type=submit] {
				width: 30vw; height: 20vh;
				font-size: 300%;
			}
		</style>
	</head>
	<body>
		<form action="/" name=input method="GET">
		<input type=submit value="Pause/Play" name=pause>
		<input type=submit value="Stop" name=stop>
		<br>
		<input type=submit value="<-30s" name=backSmall>
		<input type=submit value="30s->" name=forwardSmall>
		<br>
		<input type=submit value="<-10min" name=backBig>
		<input type=submit value="10min->" name=forwardBig>
		<br>
		<input type=submit value="Volume-" name=lower>
		<input type=submit value="Volume+" name=louder>
	</body>
</html>`

func log(a interface{}) {
	if a != nil {
		fmt.Println(a)
	}
}
