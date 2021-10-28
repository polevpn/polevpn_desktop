package main

import (
	_ "embed"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
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
		out, err := core.ExecuteCommand("bash", "-c", `/usr/bin/osascript -e "do shell script \"`+dir+`/service/polevpn_service -logPath=`+logPath+` >`+logPath+`/run.log 2>&1 &\" with prompt \"PoleVPN Request System Privileges\" with administrator privileges"`)
		if err != nil {
			glog.Error("start service fail,", err.Error()+","+string(out))
			return err
		}
	} else if runtime.GOOS == "linux" {
		out, err := core.ExecuteCommand("bash", "-c", `sudo `+dir+`/service/polevpn_service -logPath=`+logPath+` >`+logPath+`/run.log 2>&1 &`)
		if err != nil {
			glog.Error("start service fail,", err.Error()+","+string(out))
			return err
		}
	} else if runtime.GOOS == "windows" {

		go func() {
			out, err := core.ExecuteCommand(dir+`\service\polevpn_service.exe`, `-logPath=`+logPath)
			if err != nil {
				glog.Error("check servie fail,", err.Error()+","+string(out))
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
