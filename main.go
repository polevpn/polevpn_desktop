package main

import (
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/polevpn/elog"
	"github.com/webview/webview"
)

func main() {

	// flag.Parse()
	// defer elog.Flush()

	//systray.Register(onReady, func() { elog.Info("exit") })

	w := webview.New(false)

	defer w.Destroy()
	w.SetTitle("xxxxxx")
	w.SetSize(400, 600, webview.HintNone)
	w.Navigate("niubit.com/app")

	elog.Info("xxxxx")
	elog.Info("xxxxxxxxxxx")

	w.Run()

}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Webview example")
	mShowLantern := systray.AddMenuItem("Show Lantern", "")
	mShowWikipedia := systray.AddMenuItem("Show Wikipedia", "")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		for {
			select {
			case <-mShowLantern.ClickedCh:
				elog.Info("xxxxx1")
			case <-mShowWikipedia.ClickedCh:
				elog.Info("xxxxx2")
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()

}
