package main

import (
	"errors"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"os"
	"path"
	"strings"
	"time"
)

type player struct {
	Files   []string
	Current struct {
		Index      int
		Filepath   string
		Name       string
		IsPlay     bool
		FileHandle *os.File
		Streamer   beep.StreamSeekCloser
	}
	IsInit bool
}

var Player player

func (p *player) PlayList(files []string) error {

	for _, f := range files {
		if !strings.HasSuffix(f, ".mp3") {
			continue
		}
		p.Files = append(p.Files, f)
	}

	if len(p.Files) < 1 {
		return errors.New("playlist is empty")
	}
	return nil
}

func (p *player) Next() (err error) {
	if p.Current.Index > len(p.Files)-2 {
		p.Current.Index = len(p.Files) - 1
	} else {
		p.Current.Index = p.Current.Index + 1
	}
	return p.Play(p.Current.Index)
}

func (p *player) Prev() (err error) {
	if p.Current.Index < 1 {
		p.Current.Index = 0
	} else {
		p.Current.Index = p.Current.Index - 1
	}
	return p.Play(p.Current.Index)
}

func (p *player) Play(index int) (err error) {

	if len(p.Files) < 1 {
		return errors.New("playlist is empty")
	} else if index > len(p.Files)-1 {
		return errors.New("playlist file not find")
	}
	p.Current.Filepath = p.Files[index]

	name := strings.Split(path.Base(p.Current.Filepath), ".")
	p.Current.Name = name[0]

	format, err := func() (beep.Format, error) {
		speaker.Lock()
		defer speaker.Unlock()
		var format beep.Format
		if p.Current.Streamer != nil {
			p.Current.Streamer.Close()
		}
		if p.Current.FileHandle != nil {
			p.Current.FileHandle.Close()
		}

		p.Current.FileHandle, err = os.Open(p.Current.Filepath)
		if err != nil {
			return format, err
		}

		p.Current.Streamer, format, err = mp3.Decode(p.Current.FileHandle)
		if err != nil {
			return format, err
		}
		p.Current.IsPlay = true
		return format, nil
	}()
	if err != nil {
		return err
	}

	if !p.IsInit {
		p.IsInit = true
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		speaker.Play(beep.Seq(ChiStreamer{Player: p}), beep.Callback(func() {
			//p.Current.IsPlay = false
			socketio.BroadcastToAll("push", &Message{
				Name:     p.Current.Name,
				Filepath: p.Current.Filepath,
				IsPlay:   p.Current.IsPlay,
				Samples:  nil,
			})
		}))
	}

	return nil
}

func (p *player) Playing(isPlay bool) {
	fmt.Println("playing", isPlay)
	p.Current.IsPlay = isPlay
}
