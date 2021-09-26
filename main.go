package main

import (
	"flag"
	"os"

	"github.com/polevpn/elog"
	"github.com/polevpn/systray"
	"github.com/polevpn/webview"
)

var WebView webview.WebView

func main() {

	flag.Parse()
	defer elog.Flush()

	WebView = webview.New(true)
	defer WebView.Destroy()

	WebView.SetSize(300, 600, webview.HintFixed)

	dir, err := os.Getwd()
	if err != nil {
		elog.Fatal("xxxx")
	}
	WebView.Navigate("file://" + dir + "/index.html")

	WebView.Bind("hello", func(msg map[string]interface{}) string {
		elog.Info("get message ", msg)
		return "sdddddd"
	})

	systray.Register(onReady, func() { elog.Info("exit") })

	WebView.Run()
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
				WebView.Show()
			case <-mHideApp.ClickedCh:
				WebView.Hide()
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()

}
