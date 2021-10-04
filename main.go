package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/polevpn/elog"
	core "github.com/polevpn/polevpn_core"
	"github.com/polevpn/systray"
	"github.com/polevpn/webview"
)

var mainView webview.WebView
var controller *Controller
var glog *elog.EasyLogger

type LogHandler struct {
}

func (lh *LogHandler) Write(data []byte) (int, error) {
	return len(data), AddLog(string(data))
}

func (lh *LogHandler) Flush() {
}

func main() {

	flag.Parse()

	err := InitDB("./config.db")

	if err != nil {
		glog.Fatal("init db fail,", err)
	}

	glog = elog.NewEasyLogger("INFO", true, 3, &LogHandler{})
	defer glog.Flush()
	core.SetLogger(glog)

	mainView = webview.New(true, true)
	defer mainView.Destroy()
	mainView.SetSize(300, 570, webview.HintFixed)

	dir, err := os.Getwd()
	if err != nil {
		glog.Fatal("get work directory fail,", err)
	}
	mainView.Navigate("file://" + dir + "/index.html")

	controller = NewController(mainView)

	controller.Bind()

	signalHandler()

	systray.Register(onReady, func() {
		mainView.Terminate()
		controller.StopAccessServer(ReqStopAccessServer{})
		DeleteTwoDaysAgoLogs()
		glog.Info("exit")
	})

	mainView.Run()
}

func signalHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				glog.Info("receive exit signal,exit")
				systray.Quit()
			default:
			}
		}
	}()
}

func onReady() {

	iconData, err := GetIconData("polevpn.ico")
	if err != nil {
		glog.Error(err)
		return
	}
	systray.SetTooltip("PoleVPN")
	systray.SetTemplateIcon(iconData, iconData)
	systray.SetIcon(iconData)
	mShowApp := systray.AddMenuItem("Show PoleVPN", "Show PoleVPN")
	mQuit := systray.AddMenuItem("Quit", "Quit")
	go func() {
		for {
			select {
			case <-mShowApp.ClickedCh:
				mainView.Show()
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}
