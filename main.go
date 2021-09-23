package main

import (
	"flag"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/polevpn/elog"
	"github.com/webview/webview"
)

var webView webview.WebView

func main() {

	flag.Parse()
	defer elog.Flush()

	systray.Register(onReady, func() { elog.Info("exit") })

	webView = webview.New(false)
	defer webView.Destroy()

	webView.SetTitle("Niubit")
	webView.SetSize(300, 600, webview.HintNone)
	webView.Navigate("https://www.niubitstest.com/app")
	webView.Run()
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Niubit")
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
