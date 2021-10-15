./makerc.sh
go build  -ldflags="-s -w -H windowsgui"
mkdir -p build/win
mkdir -p build/win/service
cp polevpn_desktop.exe build/win/PoleVPN.exe
cp lib/Webview2Loader.dll build/win/
cd polevpn_service
go build  -ldflags="-s -w"
cp polevpn_service ../build/win/service
rm polevpn.syso
