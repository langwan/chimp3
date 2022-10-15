package main

import (
	"bytes"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"github.com/langwan/langgo/core/log"
	"io"
	"path/filepath"
	"time"

	"os"
	"strings"
)

type Streamer struct {
}

type Player struct {
	otoContext *oto.Context
	otoPlayer  *oto.Player
	Current    struct {
		Name       string
		Filepath   string
		Mp3Decoder *mp3.Decoder
		FileBuffer []byte
	}
	IsPlay        bool
	Change        chan string
	Buffer        []byte
	UpdateSamples func(p *Player, samples [][2]float64)
}

func New() *Player {

	p := Player{}
	p.otoContext, _ = oto.NewContext(44100, 2, 2, 17640)
	p.otoPlayer = p.otoContext.NewPlayer()

	p.Change = make(chan string)
	go p.Update()
	return &p
}

func (p *Player) Play(filename string) (err error) {
	p.Change <- filename
	return nil
}

func (p *Player) Update() {
	for {
		select {
		case filename := <-p.Change:
			name := strings.Split(filepath.Base(filename), ".")
			p.Current.Name = name[0]
			p.Current.Filepath = filename
			p.IsPlay = false
			var err error
			p.Current.FileBuffer, err = os.ReadFile(filename)
			if err != nil {
				continue
			}
			reader := bytes.NewReader(p.Current.FileBuffer)
			p.Current.Mp3Decoder, err = mp3.NewDecoder(reader)
			if err != nil {
				continue
			}
			p.Buffer = make([]byte, 17640)
			p.Current.Mp3Decoder.Seek(0, 0)
			p.IsPlay = true
		default:
			if p.IsPlay {
				log.Logger("backend", "player.Update").Debug().Bool("isplay", p.IsPlay).Str("name", p.Current.Name).Send()
				buf := make([]byte, 17640)
				reads := 0
				isEof := false

				for {
					read, err := p.Current.Mp3Decoder.Read(p.Buffer[:len(buf)-reads])
					for i := 0; i < read; i++ {
						buf[reads+i] = p.Buffer[i]
					}
					reads += read
					if err == io.EOF {
						isEof = true
						break
					} else if reads == len(buf) {
						break
					}
				}

				samples, err := readSamples(buf[:reads])
				if err != nil {

					return
				}
				p.UpdateSamples(p, samples)
				_, err = p.otoPlayer.Write(buf[:reads])
				if err != nil {

					return
				}

				if isEof {
					p.IsPlay = false
					p.UpdateSamples(p, nil)
				}
			} else {
				p.UpdateSamples(p, nil)
				time.Sleep(time.Second / 10)
			}
		}
	}
}

func readSamples(buf []byte) ([][2]float64, error) {
	format := Format{
		SampleRate:  44100,
		NumChannels: 2,
		Precision:   2,
	}
	samples := make([][2]float64, len(buf)/(format.NumChannels*format.Precision))
	var tmp [4]byte
	reader := bytes.NewReader(buf)
	i := 0
	for {
		_, err := reader.Read(tmp[:])
		if err == io.EOF {
			break
		} else if err != nil {
			return samples, err
		} else {
			samples[i], _ = format.DecodeSigned(tmp[:])
			i++
		}
	}
	return samples, nil
}
