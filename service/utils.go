package main

import (
	"crypto/md5"
	_ "embed"
	"encoding/hex"
	"github.com/denisbrodbeck/machineid"
	"net"
	"os"
	"strings"
	"time"
)

//go:embed version
var VERSION string

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

func GetRouteIpsFromDomain(domains []string) []string {

	ips := make([]string, 0)
	for _, domain := range domains {

		netips, err := net.LookupIP(domain)
		if err != nil {
			continue
		}
		for _, netip := range netips {
			if !strings.Contains(netip.String(), ":") {
				ips = append(ips, netip.String()+"/32")
			}
		}
	}
	return ips
}

func GetDeviceId() string {
	id, err := machineid.ID()
	if err != nil {
		return "11111111111111111111111111111111"
	}

	h := md5.New()
	h.Write([]byte(id))
	result := hex.EncodeToString(h.Sum(nil))

	return result
}
