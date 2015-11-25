package main

import (
	"io"
	"os/exec"
)

type videoPlayer struct {
	running bool
	play    *exec.Cmd
	out     io.ReadCloser
	in      io.WriteCloser
}

func (v *videoPlayer) playVideo(path string) (err error) {
	err = v.stopVideo()
	if err != nil {
		return
	}
	v.play = exec.Command("omxplayer", "-wr", path)
	v.out, err = v.play.StdoutPipe()
	if err != nil {
		return
	}
	v.in, err = v.play.StdinPipe()
	if err != nil {
		return
	}
	err = v.play.Start()
	if err != nil {
		return
	}
	v.running = true
	return nil
}

func (v *videoPlayer) stopVideo() error {
	err := v.writeIfRunning("q")
	v.running = false
	return err
}

func (v *videoPlayer) playPause() error {
	return v.writeIfRunning(" ")
}

func (v *videoPlayer) volumeDown() error {
	return v.writeIfRunning("-")
}

func (v *videoPlayer) volumeUp() error {
	return v.writeIfRunning("+")
}

func (v *videoPlayer) back30Seconds() error {
	return v.writeIfRunning(string([]byte{27, 91, 68}))
}

func (v *videoPlayer) forward30Seconds() error {
	return v.writeIfRunning(string([]byte{27, 91, 67}))
}

func (v *videoPlayer) back10Minutes() error {
	return v.writeIfRunning(string([]byte{27, 91, 66}))
}

func (v *videoPlayer) forward10Minutes() error {
	return v.writeIfRunning(string([]byte{27, 91, 65}))
}

func (v *videoPlayer) writeIfRunning(msg string) (err error) {
	if v.running {
		_, err = v.in.Write([]byte(msg))
	}
	return
}
