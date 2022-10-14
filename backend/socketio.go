package main

import (
	gosocketio "github.com/ambelovsky/gosf-socketio"
	"github.com/ambelovsky/gosf-socketio/transport"
	"github.com/gin-gonic/gin"
)

var socketio *gosocketio.Server

func NewSocketIO(g *gin.Engine) {
	socketio = gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
	g.Any("/socket.io/*any", gin.WrapH(socketio))
	socketio.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		c.Emit("hello", "im ss")
	})
}
