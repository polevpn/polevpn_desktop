go build  -ldflags="-s -w"
mkdir -p build/macos
cp -r resource/PoleVPN.app build/macos/ 
mkdir -p build/macos/PoleVPN.app/Contents/MacOS/polevpn_service
cp polevpn_desktop build/macos/PoleVPN.app/Contents/MacOS/PoleVPN
cd polevpn_service
go build  -ldflags="-s -w"
cp polevpn_service ../build/macos/PoleVPN.app/Contents/MacOS/polevpn_service/


