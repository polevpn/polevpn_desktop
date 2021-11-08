go mod tidy
go build  -ldflags="-s -w"
mkdir -p build/linux/PoleVPN
mkdir -p build/linux/PoleVPN/service
cp polevpn_desktop build/linux/PoleVPN/PoleVPN
cp version service/
cd service
go mod tidy
go build  -ldflags="-s -w"
cp -f polevpn_service ../build/linux/PoleVPN/service
cd ../build/linux/
zip -r PoleVPN.zip ./* 