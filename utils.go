package main

import (
	"io/ioutil"
	"time"
)

func GetIconData(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func getTimeTwoDaysAgo() string {
	return time.Now().Add(time.Hour * 48).Format("2006-01-02 15:04:05")
}
