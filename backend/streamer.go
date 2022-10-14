package main

import (
	"github.com/faiface/beep"
	"github.com/mjibson/go-dsp/fft"
	"math"
)

type ChiStreamer struct {
	Streamer beep.Streamer
}

func (c ChiStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	_, ok = c.Streamer.Stream(samples)
	if !ok {
		return 0, ok
	}

	//log.Logger("backend", "Stream").Debug().Int("len", len(samples)).Interface("samples", samples).Send()
	var ware = make([]float64, len(samples))
	for i := 0; i < len(samples); i++ {
		ware[i] = samples[i][0] + samples[i][1]
	}
	fftOutput = fft.FFTReal(ware)
	freqSpectrum = make([]float64, len(samples))
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
	socketio.BroadcastToAll("push", freqSpectrum)
	//window.Apply(ware, window.Blackman)
	//fftOutput = fft.FFTReal(ware)
	//updateSpectrumValues(60)
	if !isPlayer {
		isPlayer = true
	}
	return len(samples), true
}

func RangeConvert(value float64, s1 float64, s2 float64, d1 float64, d2 float64) float64 {
	w1 := s2 - s1
	w2 := d2 - d1
	return (value+s1)/w1*w2 + d1
}

func (c ChiStreamer) Err() error {
	return c.Streamer.Err()
}

//func updateSpectrumValues(maxValue float64) {
//	for i := 0; i < spectrumSize; i++ {
//		fr := real(fftOutput[i])
//		fi := imag(fftOutput[i])
//		magnitude := math.Sqrt(fr*fr + fi*fi)
//		val := math.Min(maxValue, math.Abs(magnitude))
//		if freqSpectrum[i] > val {
//			freqSpectrum[i] = math.Max(freqSpectrum[i]-peakFalloff, 0.0)
//		} else {
//			freqSpectrum[i] = (val + freqSpectrum[i]) / 2.0
//		}
//	}
//	socketio.BroadcastToAll("push", freqSpectrum)
//}
