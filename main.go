package main

import (
	"flag"
	"os"

	"github.com/getlantern/systray"
	"github.com/polevpn/elog"
	"github.com/polevpn/webview"
)

var webView webview.WebView

func main() {

	flag.Parse()
	defer elog.Flush()

	webView = webview.New(true)
	defer webView.Destroy()

	webView.SetTitle("PoleVPN")
	webView.SetSize(300, 600, webview.HintFixed)

	dir, err := os.Getwd()
	if err != nil {
		elog.Fatal("xxxx")
	}
	webView.Navigate("file://" + dir + "/index.html")

	webView.Bind("hello", func(msg map[string]interface{}) string {
		elog.Info("get message ", msg)
		return "sdddddd"
	})

	systray.Register(onReady, func() { elog.Info("exit") })

	webView.Run()
}

func onReady() {
	//systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("PoleVPN")
	mShowApp := systray.AddMenuItem("Open App", "")
	mHideApp := systray.AddMenuItem("Hide App", "")
	mQuit := systray.AddMenuItem("Quit", "Quit")
	go func() {
		for {
			select {
			case <-mShowApp.ClickedCh:
				webView.Show()
			case <-mHideApp.ClickedCh:
				webView.Hide()
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()

}
