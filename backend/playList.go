package main

import (
	"errors"
	"strings"
)

type _PlayList struct {
	Files        []string
	Player       *Player
	CurrentIndex int
}

var PlayerList _PlayList

func (p *_PlayList) PlayList(files []string) error {
	//speaker.Lock()
	//defer speaker.Unlock()
	//if p.Current.Streamer != nil {
	//	p.Current.Streamer.Close()
	//}
	//if p.Current.FileHandle != nil {
	//	p.Current.FileHandle.Close()
	//}
	p.Files = nil
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

func (p *_PlayList) Next() (err error) {
	if p.CurrentIndex > len(p.Files)-2 {
		p.CurrentIndex = len(p.Files) - 1
	} else {
		p.CurrentIndex = p.CurrentIndex + 1
	}
	return p.Play(p.CurrentIndex)
}

func (p *_PlayList) Prev() (err error) {
	if p.CurrentIndex < 1 {
		p.CurrentIndex = 0
	} else {
		p.CurrentIndex = p.CurrentIndex - 1
	}
	return p.Play(p.CurrentIndex)
}

func (p *_PlayList) Play(index int) (err error) {
	p.Player.Play(p.Files[index])
	return nil
}

func (p *_PlayList) Playing(isPlay bool) {
	p.Player.IsPlay = isPlay
}
