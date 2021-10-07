package main

import (
	_ "embed"
	"os"
	"time"
)

//go:embed resource/polevpn.ico
var iconByte []byte

func getTimeTwoHoursAgo() string {
	return time.Now().Add(time.Hour * -2).Format("2006-01-02 15:04:05")
}

func fileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
