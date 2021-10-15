package main

import (
	"io"
	"strings"

	"github.com/gorilla/websocket"
	core "github.com/polevpn/polevpn_core"
)

const (
	CH_WEBSOCKET_WRITE_SIZE = 2000
	TRAFFIC_LIMIT_INTERVAL  = 10
)

type WebSocketConn struct {
	conn    *websocket.Conn
	wch     chan []byte
	closed  bool
	handler ProcessHandler
}

func NewWebSocketConn(conn *websocket.Conn, handler ProcessHandler) *WebSocketConn {
	return &WebSocketConn{
		conn:    conn,
		closed:  false,
		wch:     make(chan []byte, CH_WEBSOCKET_WRITE_SIZE),
		handler: handler,
	}
}

func (wsc *WebSocketConn) Close(flag bool) error {
	if !wsc.closed {
		wsc.closed = true
		if wsc.wch != nil {
			wsc.wch <- nil
			close(wsc.wch)
		}
		err := wsc.conn.Close()
		if flag {
			go wsc.handler.OnClosed(wsc, false)
		}
		return err
	}
	return nil
}

func (wsc *WebSocketConn) String() string {
	return wsc.conn.RemoteAddr().String() + "->" + wsc.conn.LocalAddr().String()
}

func (wsc *WebSocketConn) IsClosed() bool {
	return wsc.closed
}

func (wsc *WebSocketConn) Read() {
	defer func() {
		wsc.Close(true)
	}()

	defer core.PanicHandler()

	for {
		mtype, pkt, err := wsc.conn.ReadMessage()
		if err != nil {
			if err == io.ErrUnexpectedEOF || err == io.EOF || strings.Contains(err.Error(), "close") {
				glog.Info(wsc.String(), ",conn closed")
			} else {
				glog.Error(wsc.String(), ",conn read exception:", err)
			}
			return
		}
		if mtype == websocket.BinaryMessage {
			wsc.handler.OnRequest(pkt, wsc)
		} else {
			glog.Info("ws mtype=", mtype)
		}
	}

}

func (wsc *WebSocketConn) Write() {
	defer core.PanicHandler()

	for {

		pkt, ok := <-wsc.wch
		if !ok {
			glog.Error(wsc.String(), ",get pkt from write channel fail,maybe channel closed")
			return
		}
		if pkt == nil {
			glog.Info(wsc.String(), ",exit write process")
			return
		}
		err := wsc.conn.WriteMessage(websocket.BinaryMessage, pkt)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF || strings.Contains(err.Error(), "close") {
				glog.Info(wsc.String(), ",conn closed")
			} else {
				glog.Error(wsc.String(), ",conn write exception:", err)
			}
			return
		}
	}
}

func (wsc *WebSocketConn) Send(pkt []byte) {
	if wsc.closed {
		return
	}
	if wsc.wch != nil {
		wsc.wch <- pkt
	}
}
