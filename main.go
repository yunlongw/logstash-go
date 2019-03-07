package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"runtime"
)

func httpLogServer() {
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

	log.Println(">>> HTTP - POST至 http://ip:9090/write")
	//开启 9090 端口监听
	http.ListenAndServe(":9090", nil)
}

func udpLogServer() {
	addr, err := net.ResolveUDPAddr("udp4", ":9091")
	if err != nil {
		log.Printf("net.ResolveUDPAddr error %q", err)
	}

	l, err := net.ListenUDP("udp4", addr)
	if err != nil {
		log.Printf("net.ListenUDP error %q", err)
	}

	defer l.Close()
	log.Println(">>> UDP - 转发至 upd://ip:9091")
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

func main() {
	log.Println("LogStation started (v1.0.0) - ws://ip:9191")
	runtime.GOMAXPROCS(runtime.NumCPU())

	//开始接收udp日志
	go udpLogServer()

	//开始接收HTTP日志,收到消息后发送广播
	go httpLogServer()

	//处理ws
	go h.run()

	//通过ws发送给客户端
	http.ListenAndServe(":9191", http.HandlerFunc(wsHandler))
}
