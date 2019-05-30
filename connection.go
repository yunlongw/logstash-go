package main

import (
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net"
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

func wbServerHandler(w http.ResponseWriter, r *http.Request)  {
	log.Println(">>> Websocket - 至 http://ip:9393/write")

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		// handle error
	}

	defer conn.Close()

	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			// handle error
			//fmt.Println(1)
		}
		err = wsutil.WriteServerMessage(conn, op, msg)
		if err != nil {
			// handle error
			//fmt.Println(2)
		}

		h.broadcastWeb <- msg
	}

	//bodyByte, _ := ioutil.ReadAll(r.Body)
	//h.broadcast <- bodyByte

}