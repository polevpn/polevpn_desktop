package main

type Conn interface {
	Send([]byte)
	Close(bool) error
	IsClosed() bool
	String() string
	Read()
	Write()
}
