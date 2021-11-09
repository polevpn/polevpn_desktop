package main

import (
	"embed"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/polevpn/elog"
	"github.com/polevpn/systray"
	"github.com/polevpn/webview"
)

var mainView webview.WebView
var glog *elog.EasyLogger

//go:embed static
var staticFiles embed.FS

func main() {

	flag.Parse()
	defer elog.Flush()
	glog = elog.GetLogger()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		glog.Fatal(err)
	}

	if !fileExist(homeDir + "/.polevpn/") {
		os.Mkdir(homeDir+"/.polevpn/", 0755)
	}

	glog.SetLogPath(homeDir + string(os.PathSeparator) + ".polevpn")

	glog.Info("check service")
	exist := CheckServiceExist()

	if !exist {
		glog.Info("starting service")
		err := StartService(homeDir + string(os.PathSeparator) + ".polevpn")
		if err != nil {
			glog.Fatal("start service fail,", err)
		}
		glog.Info("start service ok")

	}

	http.Handle("/", http.FileServer(http.FS(staticFiles)))

	go func() {
		err := http.ListenAndServe("127.0.0.1:35972", nil)
		if err != nil {
			glog.Fatal(err)
		}
	}()

	glog.Info("init db")

	err = InitDB(homeDir + "/.polevpn/config.db")

	if err != nil {
		glog.Fatal("init db fail,", err)
	}

	glog.Info("init db ok")

	mainView = webview.New(300, 570, true, true)
	defer mainView.Destroy()
	mainView.SetSize(300, 570, webview.HintFixed)
	mainView.SetIcon(iconByte)

	glog.Info("show homepage")

	mainView.Navigate("http://127.0.0.1:35972/static/index.html")

	glog.Info("init controller")
	controller, err := NewController(mainView)

	if err != nil {
		glog.Fatal("new controller fail,", err)
	}

	glog.Info("init controller ok")

	controller.Bind()

	signalHandler()
	glog.Info("register systray")
	systray.Register(onReady, func() {
		mainView.Terminate()
		controller.StopAccessServer(ReqStopAccessServer{})
		glog.Info("exit")
		glog.Flush()
	})

	glog.Info("app running")
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

	systray.SetTooltip("PoleVPN")
	systray.SetTemplateIcon(iconByte, iconByte)
	systray.SetIcon(iconByte)
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
