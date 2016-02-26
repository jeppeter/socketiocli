package socketiocli

import (
	"fmt"
	//logging "github.com/jeppeter/go-logging"
)

type SocketIoConfig struct {
	version   string
	transport string
	t         string
	base64    string
}

func newSocketIoConfig() *SocketIoConfig {
	p := &SocketIoConfig{}
	p.version = "3"
	p.transport = "polling"
	p.t = getRandName(7)
	p.base64 = "1"
	return p
}

func (cfg *SocketIoConfig) FormatQuery() string {
	return fmt.Sprintf("EIO=%s&transport=%s&t=%s&b64=%s", cfg.version,
		cfg.transport, cfg.t, cfg.base64)
}

type WSConfig struct {
	version   string
	transport string
	sid       string
}

func newWSConfig(sid string) *WSConfig {
	p := &WSConfig{}
	p.sid = sid
	p.version = "3"
	p.transport = "websocket"
	return p
}

func (ws *WSConfig) formatQuery() string {
	return fmt.Sprintf("EIO=%s&transport=%s&sid=%s", ws.version, ws.transport, ws.sid)
}
