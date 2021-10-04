module polevpn_desktop

go 1.16

require (
	github.com/polevpn/anyvalue v1.0.6 // indirect
	github.com/polevpn/elog v1.1.0
	github.com/polevpn/polevpn_core v1.0.13
	github.com/polevpn/systray v1.1.1
	github.com/polevpn/webview v0.1.2
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	gorm.io/driver/sqlite v1.1.5
	gorm.io/gorm v1.21.15
)

replace (
	github.com/polevpn/polevpn_core => ../polevpn_core
	github.com/polevpn/systray => ../systray
	github.com/polevpn/webview => ../webview
)
