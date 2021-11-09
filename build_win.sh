./makerc.sh
go mod tidy
go build  -ldflags="-s -w -H windowsgui"
rm polevpn.syso
mkdir -p build/win/PoleVPN
mkdir -p build/win/PoleVPN/service
cp polevpn_desktop.exe build/win/PoleVPN/PoleVPN.exe
cp lib/Webview2Loader.dll build/win/PoleVPN
cp -r resource/tap-windows build/win/PoleVPN
cp version service/
cd service
go mod tidy
go build  -ldflags="-s -w"
cp polevpn_service.exe ../build/win/PoleVPN/service
cd ../build/win
tar -a -c -f PoleVPN.zip ./*
