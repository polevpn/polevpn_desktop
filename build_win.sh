./makerc.sh
go build  -ldflags="-s -w -H windowsgui"
mkdir -p build/win
mkdir -p build/win/service
cp polevpn_desktop.exe build/win/PoleVPN.exe
cp lib/Webview2Loader.dll build/win/
cd service
go build  -ldflags="-s -w"
cp polevpn_service ../build/win/service
cd ../
rm polevpn.syso
