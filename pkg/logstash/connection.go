package logstash

import (
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
	"logstash/pkg/logging"
	"net"
	"net/http"
	"strings"
)

type Connection struct {
	ws   net.Conn
	send chan []byte
}

/**
发送消息
*/
func (c *Connection) write() {
	for message := range c.send {
		err := wsutil.WriteServerMessage(c.ws, ws.OpText, message)
		if err != nil {
			log.Printf("发送消息到 %s 出错 %q", c.ws.RemoteAddr(), err)
			return
		}
	}
}


func WsHandler(w http.ResponseWriter, r *http.Request) {

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Printf("UpgradeHTTP %q", err)
		return
	}

	c := &Connection{
		ws: conn,
		send: make(chan []byte, 20480),
	}

	GetHub().Register <- c

	defer func() {
		GetHub().Unregister <- c
	}()

	c.write()
}

/**
websocket 客户端接收参数
 */
func WbServerHandler(w http.ResponseWriter, r *http.Request) {

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		logging.Error(err)
	}

	defer conn.Close()

	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			logging.Error(err)
		}
		err = wsutil.WriteServerMessage(conn, op, msg)
		if err != nil {
			logging.Error(err)
		}

		if msg != nil{
			inStr := strings.TrimSpace(string(msg))
			inputs := strings.Split(inStr, " ")
			switch inputs[0] {
			case "ping":
				fmt.Println(11)
				break
			case "echo":
				fmt.Println(22)
				break
			case "quit":
				fmt.Println(33)
				break
			default:
				GetHub().BroadcastWeb <- msg
			}
		}
	}


}

func ConnHandler(c net.Conn) {
	if c == nil {
		return
	}
	buf := make([]byte, 4096)
	for {
		cnt, err := c.Read(buf)
		if err != nil || cnt == 0 {
			c.Close()
			break
		}
		inStr := strings.TrimSpace(string(buf[0:cnt]))
		inputs := strings.Split(inStr, " ")
		switch inputs[0] {
		case "ping":
			c.Write([]byte("pong\n"))
		case "echo":
			echoStr := strings.Join(inputs[1:], " ") + "\n"
			c.Write([]byte(echoStr))
		case "quit":
			c.Close()
			break
		default:
			fmt.Printf("Unsupported command: %s\n", inputs[0])
		}
	}
	fmt.Printf("Connection from %v closed. \n", c.RemoteAddr())
}
