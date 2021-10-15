package main

type ProcessHandler interface {
	OnConnected(conn Conn)
	OnRequest(pkg []byte, conn Conn)
	OnClosed(conn Conn, proactive bool)
}
