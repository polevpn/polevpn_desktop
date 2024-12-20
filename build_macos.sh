go mod tidy
go build  -ldflags="-s -w -extldflags=-mmacosx-version-min=11.0" 
mkdir -p build/macos
cp -r resource/PoleVPN.app build/macos/ 
mkdir -p build/macos/PoleVPN.app/Contents/MacOS/service
cp polevpn_desktop build/macos/PoleVPN.app/Contents/MacOS/PoleVPN
cp version service/
cd service
go mod tidy
go build  -ldflags="-s -w -extldflags=-mmacosx-version-min=11.0"
cp polevpn_service ../build/macos/PoleVPN.app/Contents/MacOS/service/
cd ../build/macos
codesign --force --deep --sign - PoleVPN.app
zip -r PoleVPN.zip ./* 


