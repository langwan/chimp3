package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/langwan/langgo"
	helperGin "github.com/langwan/langgo/helpers/gin"
	helperGrpc "github.com/langwan/langgo/helpers/grpc"
	"io"
	"sync"
	"time"
)

const peakFalloff = 8.0
const defaultWindowWidth = 800
const defaultWindowHeight = 100

var fftOutput []complex128
var fftOutputLock sync.RWMutex
var isDropped = false
var done = make(chan bool)
var isPlayer = false
var freqSpectrum []float64

func main() {
	langgo.Run()
	gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	g.Use(cors.New(cors.Config{
		AllowAllOrigins:        true,
		AllowMethods:           []string{"POST"},
		AllowHeaders:           []string{"*	"},
		AllowCredentials:       false,
		ExposeHeaders:          nil,
		MaxAge:                 12 * time.Hour,
		AllowWildcard:          false,
		AllowBrowserExtensions: false,
		AllowWebSockets:        false,
		AllowFiles:             false,
	}))
	NewSocketIO(g)
	backend.UpdateList(context.Background(), &BackendRequest{Paths: []string{"/Users/langwan/Documents/data/github/music-dl/吻别.mp3", `/Users/langwan/Documents/data/github/music-dl/房东的猫 - 海.mp3`}})
	rg := g.Group("rpc")
	rg.Any("/*uri", rpc())
	g.Run(":8000")
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

		response, code, err := helperGrpc.Call(&backend, methodName, string(body), nil)

		if err != nil {
			c.AbortWithStatus(500)
			return
		} else if code != 0 {

			helperGin.SendBad(c, code, err.Error(), nil)
		}

		helperGin.SendOk(c, response)
	}
}
