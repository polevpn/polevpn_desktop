./makerc.sh
go mod tidy
go build  -ldflags="-s -w -H windowsgui"
ldd polevpn_desktop.exe
rm polevpn.syso
mkdir -p build/win/PoleVPN
mkdir -p build/win/PoleVPN/service
mkdir -p build/win/PoleVPN/deps
cp polevpn_desktop.exe build/win/PoleVPN/PoleVPN.exe
cp resource/polevpn.ico build/win/PoleVPN/polevpn.ico
cp resource/polevpn.nsi build/win/PoleVPN/polevpn.nsi
cp resource/license.txt build/win/PoleVPN/license.txt
cp lib/Webview2Loader.dll build/win/PoleVPN
cp -r resource/drivers build/win/PoleVPN/deps
cp -r resource/MicrosoftEdgeWebview2Setup.exe build/win/PoleVPN/deps
cp -r resource/init.bat build/win/PoleVPN
cp version service/
cd service
go mod tidy
go build  -ldflags="-s -w"
cp polevpn_service.exe ../build/win/PoleVPN/service
cd ../build/win
tar -czvf PoleVPN.tgz ./*
