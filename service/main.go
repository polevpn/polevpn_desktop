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
	signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGINT, syscall.SIGQUIT:
				glog.Infof("receive exit signal %v,exit", s)
				handler.stop()
				glog.Fatal("exit")
			case syscall.SIGTERM, syscall.SIGHUP:
				glog.Infof("ignore signal %v", s)
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
