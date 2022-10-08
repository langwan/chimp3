package main

import (
	"github.com/faiface/beep"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/window"
)

type ChiStreamer struct {
	Streamer beep.Streamer
}

func (c ChiStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	var ware = make([]float64, len(samples))
	for i := 0; i < len(samples); i++ {
		ware[i] = samples[i][0] * 100
	}
	window.Apply(ware, window.Blackman)
	fftOutputLock.Lock()
	fftOutput = fft.FFTReal(ware)
	fftOutputLock.Unlock()
	if !isPlayer {
		isPlayer = true
	}
	return c.Streamer.Stream(samples)
}

func (c ChiStreamer) Err() error {
	return c.Streamer.Err()
}
