package main

import (
	"embed"
	"flag"
	"net/http"
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

//go:embed static
var staticFiles embed.FS

func main() {

	flag.Parse()

	glog = elog.NewEasyLogger("INFO", true, 3, &LogHandler{})
	defer glog.Flush()
	core.SetLogger(glog)

	http.Handle("/", http.FileServer(http.FS(staticFiles)))

	go func() {
		err := http.ListenAndServe("127.0.0.1:35972", nil)
		if err != nil {
			glog.Fatal(err)
		}
	}()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		glog.Fatal(err)
	}

	if !fileExist(homeDir + "/.polevpn/") {
		os.Mkdir(homeDir+"/.polevpn/", 0755)
	}

	err = InitDB(homeDir + "/.polevpn/config.db")

	if err != nil {
		glog.Fatal("init db fail,", err)
	}

	mainView = webview.New(true, true)
	defer mainView.Destroy()
	mainView.SetSize(300, 570, webview.HintFixed)

	mainView.Navigate("http://127.0.0.1:35972/static/index.html")

	controller = NewController(mainView)

	controller.Bind()

	signalHandler()

	systray.Register(onReady, func() {
		mainView.Terminate()
		controller.StopAccessServer(ReqStopAccessServer{})
		DeleteTwoHoursAgoLogs()
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
