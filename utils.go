package main

import (
	_ "embed"
	"errors"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	core "github.com/polevpn/polevpn_core"
)

//go:embed resource/polevpn.ico
var iconByte []byte

//go:embed version
var VERSION string

func fileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func CheckServiceExist() bool {

	_, err := http.Get("http://127.0.0.1:35973/check?version=" + VERSION)

	if err != nil {
		glog.Error("check service fail,", err)
		return false
	}
	return true
}

func StartService(logPath string) error {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		glog.Error("get current directory fail,", err)
		return err
	}

	if runtime.GOOS == "darwin" {
		err := core.RunCommand("bash", "-c", `/usr/bin/osascript -e "do shell script \"`+dir+`/service/polevpn_service -logPath=`+logPath+` >`+logPath+`/run.log 2>&1 &\" with prompt \"PoleVPN Request System Privileges\" with administrator privileges"`)
		if err != nil {
			glog.Error("start service fail,", err.Error())
			return err
		}
	} else if runtime.GOOS == "linux" {
		err := core.RunCommand("pkexec", "bash", "-c", dir+`/service/polevpn_service -logPath=`+logPath+` >`+logPath+`/run.log 2>&1 &`)
		if err != nil {
			glog.Error("start service fail,", err.Error())
			return err
		}
	} else if runtime.GOOS == "windows" {

		go func() {
			err := core.RunCommand(dir+`\service\polevpn_service.exe`, `-logPath=`+logPath)
			if err != nil {
				glog.Error("start servie fail,", err.Error())
			}
		}()
	}

	var exist bool
	for i := 0; i < 10; i++ {
		exist = CheckServiceExist()
		if exist {
			break
		}
		time.Sleep(time.Millisecond * 200)
	}
	if !exist {
		return errors.New("start service fail")
	}
	return nil
}

func GetInterfaceList() ([]string, error) {

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	outStr := make([]string, 0)

	for _, i := range interfaces {
		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			return nil, err
		}
		addresses, err := byName.Addrs()
		if err != nil {
			return nil, err
		}
		for _, v := range addresses {
			if strings.Contains(v.String(), ":") {
				continue
			}

			if strings.Contains(v.String(), "169.") || strings.Contains(v.String(), "127.") {
				continue
			}

			_, _, err := net.ParseCIDR(v.String())
			if err != nil {
				continue
			}
			outStr = append(outStr, i.Name)
		}
	}

	if len(outStr) == 0 {
		return nil, errors.New("can not find any interface")
	}

	return outStr, nil
}

func ClearDns(device string) error {

	cmd := "netsh interface ip set dns \"" + device + "\" dhcp"
	args := strings.Split(cmd, " ")

	out, err := core.ExecuteCommand(args[0], args[1:]...)

	if err != nil {
		return errors.New(err.Error() + "," + string(out))
	}
	return nil
}

func RestoreDnsServer() error {

	devices, err := GetInterfaceList()

	if err != nil {
		return errors.New("set dns fail," + err.Error())
	}

	for _, deviceName := range devices {
		err = ClearDns(deviceName)
		if err != nil {
			return errors.New("set dns fail," + err.Error())
		}
	}
	return nil
}
