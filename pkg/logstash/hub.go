package logstash

import (
	"log"
	"logstash/pkg/logging"
)

var count = 0

type Hub struct {
	Connections    map[*Connection]bool //连接
	Broadcast      chan []byte          //广播
	Register       chan *Connection     //寄存器
	registerServer chan *Connection
	Unregister     chan *Connection //注销
	BroadcastWeb   chan []byte      //websocket 接收到的消息
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register: //注册监听
			log.Println("发现新的客户端", c.ws.RemoteAddr())
			h.Connections[c] = true
		case c := <-h.Unregister: //断开链接监听
			if _, ok := h.Connections[c]; ok {
				log.Println("客户端断开连接", c.ws.RemoteAddr())
				delete(h.Connections, c)
				close(c.send)
			}
		case m := <-h.Broadcast: //消息监听
			//统计发送的数据数量
			count = count + 1
			//计算客户端数量
			clientNum := len(h.Connections)
			log.Printf("广播第 %d 条消息, 共有 %d 个客户端", count, clientNum)
			s := string(m)
			log.Printf("--- %s ----", s)
			logging.Info(s)

			for c := range h.Connections {
				select {
				case c.send <- m:
				default:
					delete(h.Connections, c)
					close(c.send)
				}
			}

		case wm := <-h.BroadcastWeb:
			//统计发送的数据数量
			count = count + 1
			//计算客户端数量
			clientNum := len(h.Connections)
			log.Printf("广播第 %d 条消息, 共有 %d 个客户端", count, clientNum)
			log.Println("---", string(wm), "---")
			logging.Info(string(wm))

			for c := range h.Connections {
				select {
				case c.send <- wm:
				default:
					delete(h.Connections, c)
					close(c.send)
				}
			}
		}
	}

}
