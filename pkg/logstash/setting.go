package logstash

import (
	"github.com/go-ini/ini"
	"log"
)

type PortList struct {
	HttpPort   string
	UdpPort    string
	WebSocket  string
	SocketPort string
}


var cfg *ini.File
var WebSetting = &PortList{}

func ConfigInit()  {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	err = cfg.Section("app").MapTo(&WebSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", "app", err)
	}
}

