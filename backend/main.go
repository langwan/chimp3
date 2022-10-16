package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/langwan/langgo"
	"github.com/langwan/langgo/core"
	"github.com/langwan/langgo/helpers/code"
	helperGin "github.com/langwan/langgo/helpers/gin"
	"github.com/mjibson/go-dsp/fft"
	"io"
	"math"
	"net/http"
	"os"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8000, "http port")
	flag.Parse()
	langgo.Run()
	gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	g.Use(cors.Default())
	NewSocketIO(g)
	PlayerList = _PlayList{Player: New()}
	PlayerList.Player.PlayList = &PlayerList
	PlayerList.Player.UpdateSamples = func(p *Player, samples [][2]float64) {
		if samples == nil {
			if socketio != nil {
				socketio.BroadcastToAll("push", &struct {
					Name     string    `json:"name"`
					Filepath string    `json:"filepath"`
					IsPlay   bool      `json:"is_player"`
					Samples  []float64 `json:"samples"`
					Mode     int       `json:"mode"`
				}{
					Name:     p.Current.Name,
					Filepath: p.Current.Filepath,
					IsPlay:   false,
					Samples:  nil,
					Mode:     p.Mode,
				})
			}

			return
		}
		var ware [2][]float64
		ware[0] = make([]float64, len(samples))
		ware[1] = make([]float64, len(samples))
		for i := 0; i < len(samples); i++ {
			ware[0][i] = samples[i][0]
			ware[1][i] = samples[i][1]
		}

		wareReal := fft.FFTReal(ware[0])
		var max float64 = 0
		for i := 0; i < len(samples); i++ {
			fr := real(wareReal[i])
			fi := imag(wareReal[i])
			magnitude := math.Sqrt(fr*fr + fi*fi)
			ware[0][i] = magnitude
			if magnitude > max {
				max = magnitude
			}
		}
		for i := 0; i < len(samples); i++ {
			ware[0][i] = RangeConvert(ware[0][i], 0, max, 0, 60)
		}
		if p.Mode == 0 || p.Mode == 3 {
			wareReal = fft.FFTReal(ware[1])
		} else {
			wareReal = fft.FFTReal(ware[0])
		}

		max = 0
		for i := 0; i < len(samples); i++ {
			fr := real(wareReal[i])
			fi := imag(wareReal[i])
			magnitude := math.Sqrt(fr*fr + fi*fi)
			ware[1][i] = magnitude
			if magnitude > max {
				max = magnitude
			}
		}
		for i := 0; i < len(samples); i++ {
			ware[1][i] = RangeConvert(ware[1][i], 0, max, 0, 60)
		}

		if socketio != nil {
			socketio.BroadcastToAll("push", &struct {
				Name     string       `json:"name"`
				Filepath string       `json:"filepath"`
				IsPlay   bool         `json:"is_player"`
				Samples  [2][]float64 `json:"samples"`
				Mode     int          `json:"mode"`
			}{
				Name:     p.Current.Name,
				Filepath: p.Current.Filepath,
				IsPlay:   p.IsPlay,
				Samples:  ware,
				Mode:     p.Mode,
			})
		}
	}
	rg := g.Group("/rpc")
	rg.Any("/*uri", rpc())

	if core.EnvName == core.Production {
		g.StaticFS("app", http.Dir("./frontend"))
		g.NoRoute(func(c *gin.Context) {
			c.File("./frontend/index.html")
		})
	}
	host := fmt.Sprintf(":%d", port)
	fmt.Printf("http start %s\n", host)
	g.Run(host)
}

func rpc() gin.HandlerFunc {

	return func(c *gin.Context) {

		methodName := c.Param("uri")[1:]

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return
		}

		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		if methodName == "Quit" {
			helperGin.SendOk(c, "ok")
			os.Exit(0)
		}

		response, code, err := code.Call(context.Background(), &backend, methodName, string(body))

		if err != nil {
			c.AbortWithStatus(500)
			return
		} else if code != 0 {

			helperGin.SendBad(c, code, err.Error(), nil)
		}

		helperGin.SendOk(c, response)
	}
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.Status(204)
			return
		}

		c.Next()
	}
}
