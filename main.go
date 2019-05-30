package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"runtime"
)

func httpLogServer(port int) {
	http.HandleFunc("/write", func(w http.ResponseWriter, r *http.Request) {

		//允许跨域访问
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Action, Module")
		}

		bodyByte, _ := ioutil.ReadAll(r.Body)
		h.broadcast <- bodyByte
	})

	//开启 9090 端口监听
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func udpLogServer(port int) {
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("net.ResolveUDPAddr error %q", err)
	}

	l, err := net.ListenUDP("udp4", addr)
	if err != nil {
		log.Printf("net.ListenUDP error %q", err)
	}

	defer l.Close()

	for {
		buf := make([]byte, 40960)
		length, _, err := l.ReadFrom(buf)
		if err != nil {
			log.Printf("l.ReadFrom err %q", err)
		}

		if length > 0 {
			h.broadcast <- buf[:length]
		} else {
			continue
		}
	}
}

func socket(port int) {
	server, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		fmt.Printf("Fail to start server, %s\n", err)
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Printf("Fail to connect, %s\n", err)
			break
		}

		go connHandler(conn)
	}
}

func main() {
	log.Println("LogStation started (v1.0.0) - ws://ip:9191")

	runtime.GOMAXPROCS(runtime.NumCPU())

	//开始接收udp日志
	go udpLogServer(9193)
	log.Println(">>> UDP - 转发至 upd://ip:9193")

	//开始接收HTTP日志
	go httpLogServer(9192)
	log.Println(">>> HTTP - POST至 http://ip:9192/write")

	// socket
	go socket(9195)

	//处理ws
	go h.run()

	//开始接受 websocket 日志
	go http.ListenAndServe(fmt.Sprintf(":%d", 9194), http.HandlerFunc(wbServerHandler))
	log.Println(">>> Websocket - 至 http://ip:9194/write")

	//通过ws发送给客户端
	http.ListenAndServe(fmt.Sprintf(":%d", 9191), http.HandlerFunc(wsHandler))
}
