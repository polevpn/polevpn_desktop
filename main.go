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

	WebView = webview.New(true, false)
	defer WebView.Destroy()
	WebView.SetSize(300, 575, webview.HintFixed)

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

	iconData, err := GetIconData("polevpn.ico")
	if err != nil {
		elog.Error(err)
		return
	}
	elog.Info(iconData)
	systray.SetTooltip("PoleVPN")
	systray.SetTemplateIcon(iconData, iconData)
	systray.SetIcon(iconData)
	mShowApp := systray.AddMenuItem("Open App", "OpenApp")
	mHideApp := systray.AddMenuItem("Hide App", "HideApp")
	mQuit := systray.AddMenuItem("Quit", "Quit")
	go func() {
		for {
			select {
			case <-mShowApp.ClickedCh:
				WebView.Show()
			case <-mHideApp.ClickedCh:
				WebView.Hide()
			case <-mQuit.ClickedCh:
				WebView.Terminate()
				systray.Quit()
			}
		}
	}()
}
