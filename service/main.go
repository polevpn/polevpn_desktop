package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/polevpn/elog"
)

var glog *elog.EasyLogger
var handler *RequestHandler

func signalHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				glog.Info("receive exit signal,exit")
				handler.stop()
				os.Exit(0)
			default:
			}
		}
	}()
}

func main() {

	flag.Parse()
	defer elog.Flush()
	glog = elog.GetLogger()

	handler = NewRequestHandler()

	server := NewHttpServer(handler)

	signalHandler()

	err := server.Listen("127.0.0.1:35973")

	if err != nil {
		glog.Fatal("start service fail,", err)
	}
}
