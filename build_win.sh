./makerc.sh
go build  -ldflags="-s -w -H windowsgui"
mkdir -p build/win
cp polevpn_desktop.exe build/win/Polevpn.exe
cp lib/Webview2Loader.dll build/win/
rm polevpn.syso
