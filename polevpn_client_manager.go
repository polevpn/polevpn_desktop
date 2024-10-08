package main

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/polevpn/anyvalue"
)

type PoleVPNClientManager struct {
	mutex    *sync.Mutex
	callback func(av *anyvalue.AnyValue)
	conn     Conn
	worked   bool
}

func NewPoleVPNClientManager(callback func(av *anyvalue.AnyValue)) (*PoleVPNClientManager, error) {

	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:35973", nil)

	if err != nil {
		return nil, err
	}

	mgr := &PoleVPNClientManager{
		mutex:    &sync.Mutex{},
		callback: callback,
		worked:   false,
	}

	mgr.conn = NewWebSocketConn(conn, mgr)

	go mgr.conn.Read()
	go mgr.conn.Write()

	return mgr, nil

}

func (pcm *PoleVPNClientManager) OnConnected(conn Conn) {

}

func (pcm *PoleVPNClientManager) OnClosed(conn Conn, proactive bool) {

	if pcm.callback != nil {
		pcm.callback(anyvalue.New().Set("event", "stoped"))
	}

}

func (pcm *PoleVPNClientManager) OnRequest(pkg []byte, conn Conn) {

	pcm.worked = true

	av, err := anyvalue.NewFromJson(pkg)
	if err != nil {
		glog.Error("decode json fail,", err)
		return
	}
	if pcm.callback != nil {
		pcm.callback(av)
	}
}

func (pcm *PoleVPNClientManager) GetLogs() {

	req := anyvalue.New()
	req.Set("event", "getlogs")
	pkt, _ := req.EncodeJson()
	pcm.conn.Send(pkt)
}

func (pcm *PoleVPNClientManager) GetUpDownBytes() {

	req := anyvalue.New()
	req.Set("event", "getbytes")
	pkt, _ := req.EncodeJson()
	pcm.conn.Send(pkt)
}

func (pcm *PoleVPNClientManager) Start(server AccessServer) error {

	if pcm.conn.IsClosed() {
		return errors.New("service stopped,please restart app")
	}

	pcm.worked = false

	req := anyvalue.New()
	req.Set("event", "start")
	req.Set("data", &server)

	pkt, _ := req.EncodeJson()
	pcm.conn.Send(pkt)

	timer := time.NewTimer(time.Second * 15)

	go func() {
		<-timer.C
		if !pcm.worked {
			//service no responses,kill service
			http.Get("http://127.0.0.1:35973/check?version=kill")
		}
	}()

	return nil
}

func (pcm *PoleVPNClientManager) Stop() error {

	req := anyvalue.New()
	req.Set("event", "stop")
	pkt, _ := req.EncodeJson()
	pcm.conn.Send(pkt)
	return nil
}
