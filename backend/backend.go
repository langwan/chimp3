package main

import (
	"errors"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/langwan/langgo/core/log"
	"os"
	"strings"
	"time"
)

type Backend struct {
}

var backend = Backend{}

type Empty struct {
}

type BackendRequest struct {
	Paths []string
}

func (b *Backend) UpdateList(request *BackendRequest) (*Empty, error) {

	if !strings.HasSuffix(request.Paths[0], ".mp3") {
		return &Empty{}, errors.New("is not mp3 file")
	}
	f, err := os.Open(request.Paths[0])
	if err != nil {
		return &Empty{}, err
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return &Empty{}, err
	}
	log.Logger("backend", "UpdateList").Debug().Interface("format", format).Send()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(beep.Seq(ChiStreamer{Streamer: streamer}, beep.Callback(func() {
		done <- true
	})))
	isDropped = true

	return &Empty{}, nil
}
