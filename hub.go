package main

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
)

var count = 0

type hub struct {
	connections    map[*connection]bool //连接
	broadcast      chan []byte          //广播
	register       chan *connection     //寄存器
	registerServer chan *connection
	unregister     chan *connection //注销
	broadcastWeb   chan []byte      //websocket 接收到的消息
}

var h = hub{
	broadcast:    make(chan []byte),
	broadcastWeb: make(chan []byte),
	register:     make(chan *connection),
	unregister:   make(chan *connection),
	connections:  make(map[*connection]bool),
}

type JsonData struct {
	Data    string `json:"data"`
	Get     string `json:"get"`
	Post    string `json:"post"`
	Cookie  string `json:"cookie"`
	Session string `json:"session"`
}

func ProtectRun(entry func()) {
	// 延迟处理的函数
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		err := recover()
		switch err.(type) {
		case runtime.Error: // 运行时错误
			fmt.Println("runtime error:", err)
		default: // 非运行时错误
			fmt.Println("error:", err)
		}
	}()
	entry()
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register: //注册监听
			log.Println("发现新的客户端", c.ws.RemoteAddr())
			h.connections[c] = true
		case c := <-h.unregister: //断开链接监听
			if _, ok := h.connections[c]; ok {
				log.Println("客户端断开连接", c.ws.RemoteAddr())
				delete(h.connections, c)
				close(c.send)
			}
		case m := <-h.broadcast: //消息监听
			//统计发送的数据数量
			count = count + 1
			//计算客户端数量
			clientNum := len(h.connections)
			log.Printf("广播第 %d 条消息, 共有 %d 个客户端", count, clientNum)
			log.Println("---", string(m), "---")

			//2019/03/07 15:54:01 invalid character ' ' in string escape code
			//2019/03/07 16:16:03 unexpected end of JSON input
			var b = false
			ProtectRun(func() {
				data := &JsonData{}
				//json 解析为 struct
				err := json.Unmarshal(m, &data)
				if err != nil {
					b = true
				}

				//转换数据
				d, err := json.Marshal(data)
				if err != nil {
					b = true
				}

				/**
				b = true
				/**
				Go 语言中 range 关键字用于 for 循环中迭代数组(array)、切片(slice)、通道(channel)或集合(map)的元素。
				在数组和切片中它返回元素的索引和索引对应的值，在集合中返回 key-value 对的 key 值。
				*/

				if b == false {
					for c := range h.connections {
						//select是Go中的一个控制结构，类似于用于通信的switch语句。每个case必须是一个通信操作，要么是发送要么是接收。
						//select随机执行一个可运行的case。如果没有case可运行，它将阻塞，直到有case可运行。一个默认的子句应该总是可运行的。
						select {
						case c.send <- d:
						default:
							delete(h.connections, c)
							close(c.send)
						}
					}
				}

			})

			if b == true {
				s := string(m)
				//default_json := strings.Join([]string{s},"")
				//default_json := fmt.Sprintf("{\"%s\":\"%s\"}","error",s)
				default_json2 := "[]"

				var data = &JsonData{Data: s, Get: default_json2, Post: default_json2, Cookie: default_json2, Session: default_json2}
				d, _ := json.Marshal(data)

				for c := range h.connections {
					//select是Go中的一个控制结构，类似于用于通信的switch语句。每个case必须是一个通信操作，要么是发送要么是接收。
					//select随机执行一个可运行的case。如果没有case可运行，它将阻塞，直到有case可运行。一个默认的子句应该总是可运行的。
					select {
					case c.send <- d:
					default:
						delete(h.connections, c)
						close(c.send)
					}
				}
			}

			println("123123123")

		case wm := <-h.broadcastWeb:
			//统计发送的数据数量
			count = count + 1
			//计算客户端数量
			clientNum := len(h.connections)
			log.Printf("广播第 %d 条消息, 共有 %d 个客户端", count, clientNum)
			log.Println("---", string(wm), "---")
			data := &JsonData{}
			//转换数据
			d, err := json.Marshal(data)
			if err != nil {
				fmt.Println(err)
			}

			for c := range h.connections {
				//select是Go中的一个控制结构，类似于用于通信的switch语句。每个case必须是一个通信操作，要么是发送要么是接收。
				//select随机执行一个可运行的case。如果没有case可运行，它将阻塞，直到有case可运行。一个默认的子句应该总是可运行的。
				select {
				case c.send <- d:
				default:
					delete(h.connections, c)
					close(c.send)
				}
			}
		}
	}

}
