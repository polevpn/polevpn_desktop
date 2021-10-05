package main

import (
	"errors"
	"runtime"
	"strings"
	"sync"

	"github.com/polevpn/anyvalue"
	core "github.com/polevpn/polevpn_core"
)

const (
	CLIENT_STOPPED  = 0
	CLIENT_STARTING = 1
	CLIENT_STARTED  = 2
)

type PoleVPNClientManager struct {
	mutex      *sync.Mutex
	status     int
	client     *core.PoleVpnClient
	server     AccessServer
	networkmgr core.NetworkManager
	device     *core.TunDevice
	callback   func(av *anyvalue.AnyValue)
}

func NewPoleVPNClientManager(callback func(av *anyvalue.AnyValue)) *PoleVPNClientManager {
	return &PoleVPNClientManager{
		mutex:    &sync.Mutex{},
		status:   CLIENT_STOPPED,
		callback: callback,
	}
}

func (pcm *PoleVPNClientManager) eventHandler(event int, client *core.PoleVpnClient, av *anyvalue.AnyValue) {

	switch event {
	case core.CLIENT_EVENT_ADDRESS_ALLOCED:
		{
			pcm.mutex.Lock()
			defer pcm.mutex.Unlock()

			pcm.status = CLIENT_STARTED

			var err error
			var routes []string

			if pcm.server.UseRemoteRouteRules {
				routes = append(routes, av.Get("route").AsStrArr()...)
			}

			if pcm.server.LocalRouteRules != "" {
				routes = append(routes, strings.Split(pcm.server.LocalRouteRules, ",")...)
			}

			glog.Info("route=", routes, ",allocated ip=", av.Get("ip").AsStr(), ",dns=", av.Get("dns").AsStr())

			if runtime.GOOS == "windows" {
				err = pcm.device.GetInterface().SetTunNetwork(av.Get("ip").AsStr() + "/30")
				if err != nil {
					glog.Error("set tun network fail,", err)
					client.Stop()
					return
				}
			}
			av.Set("remoteIp", client.GetRemoteIP())
			err = pcm.networkmgr.SetNetwork(pcm.device.GetInterface().Name(), av.Get("ip").AsStr(), client.GetRemoteIP(), av.Get("dns").AsStr(), routes)
			if err != nil {
				glog.Error("set network fail,", err)
				client.Stop()
				return
			}
			if pcm.callback != nil {
				pcm.callback(anyvalue.New().Set("event", "allocated").Set("data", av.AsMap()))
			}
		}
	case core.CLIENT_EVENT_STOPPED:
		{
			glog.Info("client stoped")
			pcm.networkmgr.RestoreNetwork()
			if pcm.callback != nil {
				pcm.callback(anyvalue.New().Set("event", "stoped").Set("data", nil))
			}
			pcm.status = CLIENT_STOPPED
		}
	case core.CLIENT_EVENT_RECONNECTED:
		glog.Info("client reconnected")
		if pcm.callback != nil {
			pcm.callback(anyvalue.New().Set("event", "reconnected").Set("data", nil))
		}
	case core.CLIENT_EVENT_RECONNECTING:
		err := pcm.networkmgr.RefreshDefaultGateway()
		if err != nil {
			glog.Error("refresh default gateway fail,", err)
		}
		glog.Info("client reconnecting")
		if pcm.callback != nil {
			pcm.callback(anyvalue.New().Set("event", "reconnecting").Set("data", nil))
		}
	case core.CLIENT_EVENT_STARTED:
		glog.Info("client started")

		var err error
		pcm.device, err = core.NewTunDevice()
		if err != nil {
			glog.Error("create device fail,", err)
			go client.Stop()
			return
		}

		client.AttachTunDevice(pcm.device)

		if pcm.callback != nil {
			pcm.callback(anyvalue.New().Set("event", "started").Set("data", nil))
		}
	case core.CLIENT_EVENT_ERROR:
		glog.Info("Unexception error,", av.Get("error").AsStr())
		if pcm.callback != nil {
			pcm.callback(anyvalue.New().Set("event", "error").Set("data", av.AsMap()))
		}
	default:
		glog.Error("invalid event=", event)
	}

}

func (pcm *PoleVPNClientManager) GetUpDownBytes() (uint64, uint64) {
	if pcm.client != nil {
		return pcm.client.GetUpDownBytes()
	}
	return 0, 0
}

func (pcm *PoleVPNClientManager) Start(server AccessServer) error {
	pcm.mutex.Lock()
	defer pcm.mutex.Unlock()

	if pcm.status != CLIENT_STOPPED {
		return errors.New("client have started")
	}

	pcm.status = CLIENT_STARTING

	glog.Info("Connect to ", server.Endpoint)

	pcm.server = server

	var err error
	pcm.client, err = core.NewPoleVpnClient()

	if err != nil {
		return err
	}

	if runtime.GOOS == "darwin" {
		pcm.networkmgr = core.NewDarwinNetworkManager()
	} else if runtime.GOOS == "linux" {
		pcm.networkmgr = core.NewLinuxNetworkManager()
	} else if runtime.GOOS == "windows" {
		pcm.networkmgr = core.NewWindowsNetworkManager()
	} else {
		return errors.New("os platform not support")
	}

	pcm.client.SetEventHandler(pcm.eventHandler)

	go pcm.client.Start(server.Endpoint, server.User, server.Password, server.Sni, server.SkipVerifySSL)

	return nil
}

func (pcm *PoleVPNClientManager) Stop() error {
	pcm.mutex.Lock()
	defer pcm.mutex.Unlock()

	if pcm.status != CLIENT_STARTED {
		return errors.New("client haven't started")
	}

	if pcm.client != nil {
		pcm.client.Stop()
	}
	pcm.status = CLIENT_STOPPED
	return nil
}
