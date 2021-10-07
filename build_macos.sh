go build  -ldflags="-s -w"
mkdir -p build/macos
cp -r resource/Polevpn.app build/macos/ 
cp polevpn_desktop build/macos/Polevpn.app/Contents/MacOS/polevpnclient.app/polevpnclient

