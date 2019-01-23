package main

import (
	"net"
	"github.com/gobwas/ws/wsutil"
	"github.com/gobwas/ws"
	"log"
	"net/http"
)

type connection struct {
	ws   net.Conn
	send chan []byte
}

func (c *connection) reader() {

}

/**
发送消息
 */
func (c *connection) write() {
	for message := range c.send {
		err := wsutil.WriteServerMessage(c.ws, ws.OpText, message)
		if err != nil {
			log.Printf("发送消息到 %s 出错 %q", c.ws.RemoteAddr(), err)
			return
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Printf("UpgradeHTTP %q", err)
		return
	}

	c := &connection{ws: conn, send: make(chan []byte, 20480)}
	h.register <- c

	defer func() {
		h.unregister <- c
	}()

	c.write()
}
