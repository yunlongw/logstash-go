package logstash

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"runtime"
)

var hub *Hub

func GetHub() *Hub {
	if hub == nil{
		hub = &Hub{
			Broadcast:    make(chan []byte),
			BroadcastWeb: make(chan []byte),
			Register:     make(chan *Connection),
			Unregister:   make(chan *Connection),
			Connections:  make(map[*Connection]bool),
		}
	}
	return hub
}

func Start() {
	ConfigInit()
	GetHub()

	//runtime.GOMAXPROCS(逻辑CPU数量)
	runtime.GOMAXPROCS(runtime.NumCPU())

	//开始接收udp日志s
	go UdpLogServer(hub)

	//开始接收HTTP日志
	go HttpLogServer(hub)

	// socket
	//go Socket()

	//处理ws
	go hub.Run()

	//开始接受 websocket 日志
	//go http.ListenAndServe(fmt.Sprintf(":%d", 9194), http.HandlerFunc(WbServerHandler))

	//通过ws发送给客户端
	log.Printf(">>>Websocket - 至 http://ip:%s/write\n", WebSetting.SocketPort)
	fmt.Printf(">>>HTTP - POST至 http://ip:%s/write\n", WebSetting.UdpPort)
	fmt.Printf(">>>UDP - 转发至 upd://ip:%s\n", WebSetting.UdpPort)
	fmt.Printf("客户端 -  ws://ip:%s\n", WebSetting.WebSocket)


	// 日志打印客户端
	http.ListenAndServe(fmt.Sprintf(":%s", WebSetting.WebSocket), http.HandlerFunc(WsHandler))
}

func HttpLogServer(h *Hub) {
	http.HandleFunc("/write", func(w http.ResponseWriter, r *http.Request) {

		//允许跨域访问
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Action, Module")
		}

		bodyByte, _ := ioutil.ReadAll(r.Body)
		h.Broadcast <- bodyByte
	})

	http.ListenAndServe(fmt.Sprintf(":%s", WebSetting.HttpPort), nil)
}

func UdpLogServer(h *Hub) {
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%s", WebSetting.UdpPort))
	if err != nil {
		log.Printf("net.ResolveUDPAddr error %q", err)
	}

	server, err := net.ListenUDP("udp4", addr)
	if err != nil {
		log.Printf("net.ListenUDP error %q", err)
	}

	defer server.Close()

	for {
		buf := make([]byte, 40960)
		length, _, err := server.ReadFrom(buf)
		if err != nil {
			log.Printf("l.ReadFrom err %q", err)
		}

		if length > 0 {
			h.Broadcast <- buf[:length]
		} else {
			continue
		}
	}
}

func Socket() {
	server, err := net.Listen("tcp",fmt.Sprintf(":%s", WebSetting.SocketPort) )

	if err != nil {
		fmt.Printf("Fail to start server, %s\n", err)
	}

	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Printf("Fail to connect, %s\n", err)
			break
		}

		go ConnHandler(conn)
	}
}
