package main

import (
	"errors"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"os"
	"strings"
	"sync"
	"time"
)

const peakFalloff = 8.0
const defaultWindowWidth = 800
const defaultWindowHeight = 100
const spectrumSize = 80

var fftOutput []complex128
var fftOutputLock sync.RWMutex
var isDropped = false
var done = make(chan bool)
var isPlayer = false

func main() {
	freqSpectrum := make([]float64, spectrumSize)
	var windowWidth int32 = defaultWindowWidth
	var windowHeight int32 = defaultWindowHeight
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(windowWidth, windowHeight, "chimp3")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	var droppedFiles []string

	for !rl.WindowShouldClose() {
		windowWidth = int32(rl.GetScreenWidth())
		windowHeight = int32(rl.GetScreenHeight())
		columnWidth := int32(windowWidth / spectrumSize)
		if rl.IsFileDropped() {
			droppedFiles = rl.LoadDroppedFiles()
			rl.UnloadDroppedFiles()
			handleFileDrop(droppedFiles[0])

			fmt.Println(droppedFiles[0])
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		if !isDropped {
			drawDropzone(windowWidth, windowHeight)
		} else {

			select {
			case <-done:
				isDropped = false
			default:
			}
			if isPlayer {
				fftOutputLock.RLock()
				updateSpectrumValues(float64(windowHeight), freqSpectrum)
				fftOutputLock.RUnlock()
				for i, s := range freqSpectrum {
					rl.DrawRectangleGradientV(int32(i)*columnWidth, windowHeight-int32(s), columnWidth, int32(s), rl.Orange, rl.Green)
					rl.DrawRectangleLines(int32(i)*columnWidth, windowHeight-int32(s), columnWidth, int32(s), rl.Black)
				}

			}

		}
		rl.EndDrawing()
	}

}

func handleFileDrop(path string) error {
	if !strings.HasSuffix(path, ".mp3") {
		return errors.New("is not mp3 file")
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return err
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(beep.Seq(ChiStreamer{Streamer: streamer}, beep.Callback(func() {
		done <- true
	})))
	isDropped = true
	return nil
}

func drawDropzone(windowWidth, windowHeight int32) {
	var fontSize float32 = 16.0
	font := rl.GetFontDefault()
	message := "Drop your mp3 to this window!"
	textPos := rl.Vector2{
		X: float32(windowWidth)/2.0 - rl.MeasureTextEx(font, message, fontSize, 2).X/2.0,
		Y: float32(windowHeight)/2.0 - fontSize/2.0}
	rl.DrawTextEx(font, message, textPos, fontSize, 2, rl.White)
	rl.DrawRectangleLines(20, 20, windowWidth-40, windowHeight-40, rl.LightGray)
}

func updateSpectrumValues(maxValue float64, freqSpectrum []float64) {
	for i := 0; i < spectrumSize; i++ {
		fr := real(fftOutput[i])
		fi := imag(fftOutput[i])
		magnitude := math.Sqrt(fr*fr + fi*fi)
		val := math.Min(maxValue, math.Abs(magnitude))
		if freqSpectrum[i] > val {
			freqSpectrum[i] = math.Max(freqSpectrum[i]-peakFalloff, 0.0)
		} else {
			freqSpectrum[i] = (val + freqSpectrum[i]) / 2.0
		}
	}
}
