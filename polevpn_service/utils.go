package main

import (
	"os"
	"strings"
	"time"
)

func GetAppName() string {
	appName := os.Args[0]
	slash := strings.LastIndex(appName, string(os.PathSeparator))
	if slash >= 0 {
		appName = appName[slash+1:]
	}
	return appName
}

func GetTimeNowDate() string {
	return time.Now().Format("2006-01-02")
}
