package main

import "io/ioutil"

func GetIconData(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
