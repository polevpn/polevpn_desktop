module polevpn_desktop

go 1.16

require (
	github.com/polevpn/elog v1.1.0
	github.com/polevpn/systray v1.1.1
	github.com/polevpn/webview v0.1.2
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
)

replace (
	github.com/polevpn/systray => ../systray
	github.com/polevpn/webview => ../webview
)
