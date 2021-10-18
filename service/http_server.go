package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	core "github.com/polevpn/polevpn_core"
)

const (
	TCP_WRITE_BUFFER_SIZE = 524288
	TCP_READ_BUFFER_SIZE  = 524288
)

type HttpServer struct {
	handler  ProcessHandler
	upgrader *websocket.Upgrader
}

func NewHttpServer(handler ProcessHandler) *HttpServer {

	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: false,
	}

	return &HttpServer{handler: handler, upgrader: upgrader}
}

func (hs *HttpServer) defaultHandler(w http.ResponseWriter, r *http.Request) {
	hs.respError(http.StatusForbidden, w)
}

func (hs *HttpServer) check(w http.ResponseWriter, r *http.Request) {
	version := r.URL.Query().Get("version")
	if version != VERSION {
		glog.Fatal("version not equal,", version, ",", VERSION)
	} else {
		w.Write([]byte("ok"))
	}
}

func (hs *HttpServer) Listen(addr string) error {

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/" {
			hs.wsHandler(w, r)
		} else if r.URL.Path == "/check" {
			hs.check(w, r)
		} else {
			hs.defaultHandler(w, r)
		}
	})
	return http.ListenAndServe(addr, handler)
}

func (hs *HttpServer) respError(status int, w http.ResponseWriter) {
	if status == http.StatusBadRequest {
		w.Header().Add("Server", "nginx/1.10.3")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<html>\n<head><title>400 Bad Request</title></head>\n<body bgcolor=\"white\">\n<center><h1>400 Bad Request</h1></center>\n<hr><center>nginx/1.10.3</center>\n</body>\n</html>"))
	} else if status == http.StatusForbidden {
		w.Header().Add("Server", "nginx/1.10.3")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("<html>\n<head><title>403 Forbidden</title></head>\n<body bgcolor=\"white\">\n<center><h1>403 Forbidden</h1></center>\n<hr><center>nginx/1.10.3</center>\n</body>\n</html>"))

	}
}

func (hs *HttpServer) wsHandler(w http.ResponseWriter, r *http.Request) {

	defer core.PanicHandler()

	conn, err := hs.upgrader.Upgrade(w, r, nil)
	if err != nil {
		glog.Error("upgrade http request to ws fail", err)
		return
	}

	glog.Info("connect attached ", conn.RemoteAddr().String())

	if hs.handler != nil {
		wsconn := NewWebSocketConn(conn, hs.handler)
		hs.handler.OnConnected(wsconn)
		go wsconn.Read()
		go wsconn.Write()
	} else {
		glog.Error("ws conn handler haven't set")
		conn.Close()
	}

}
