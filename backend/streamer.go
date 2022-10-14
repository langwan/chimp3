package main

import (
	"github.com/mjibson/go-dsp/fft"
	"math"
)

type Message struct {
	Name     string    `json:"name"`
	Filepath string    `json:"filepath"`
	IsPlay   bool      `json:"is_player"`
	Samples  []float64 `json:"samples"`
}

type ChiStreamer struct {
	Player *player
}

func (c ChiStreamer) Stream(samples [][2]float64) (n int, ok bool) {

	freqSpectrum = make([]float64, len(samples))
	if c.Player.Current.Streamer == nil {
		return 0, false
	}
	if !c.Player.Current.IsPlay {
		for i := range samples {
			samples[i] = [2]float64{}
		}

		for i := 0; i < len(samples); i++ {
			freqSpectrum[i] = 0
		}
		socketio.BroadcastToAll("push", &Message{
			Name:     c.Player.Current.Name,
			Filepath: c.Player.Current.Filepath,
			IsPlay:   c.Player.Current.IsPlay,
			Samples:  nil,
		})
		return len(samples), true
	}
	//log.Logger("backend", "Stream").Debug().Int("len", len(samples)).Interface("samples", samples).Send()
	var ware = make([]float64, len(samples))
	for i := 0; i < len(samples); i++ {
		ware[i] = samples[i][0] + samples[i][1]
	}
	fftOutput = fft.FFTReal(ware)

	//window.Apply(ware, window.Blackman)
	var max float64 = 0
	for i := 0; i < len(samples); i++ {
		fr := real(fftOutput[i])
		fi := imag(fftOutput[i])
		magnitude := math.Sqrt(fr*fr + fi*fi)
		freqSpectrum[i] = magnitude
		if magnitude > max {
			max = magnitude
		}
	}
	//log.Logger("backend", "stream").Debug().Float64("max", max).Send()
	for i := 0; i < len(samples); i++ {
		freqSpectrum[i] = RangeConvert(freqSpectrum[i], 0, max, 0, 60)
	}
	//log.Logger("backend", "stream").Debug().Interface("freqSpectrum", freqSpectrum).Send()
	socketio.BroadcastToAll("push", &Message{
		Name:     c.Player.Current.Name,
		Filepath: c.Player.Current.Filepath,
		IsPlay:   c.Player.Current.IsPlay,
		Samples:  freqSpectrum,
	})
	//window.Apply(ware, window.Blackman)
	//fftOutput = fft.FFTReal(ware)
	//updateSpectrumValues(60)
	if !isPlayer {
		isPlayer = true
	}
	return c.Player.Current.Streamer.Stream(samples)
}

func RangeConvert(value float64, s1 float64, s2 float64, d1 float64, d2 float64) float64 {
	w1 := s2 - s1
	w2 := d2 - d1
	return (value+s1)/w1*w2 + d1
}

func (c ChiStreamer) Err() error {
	return c.Player.Current.Streamer.Err()
}
