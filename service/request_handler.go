package main

import (
	"errors"
	"io/ioutil"
	"os"
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

type RequestHandler struct {
	conn       Conn
	mutex      *sync.Mutex
	status     int
	server     *anyvalue.AnyValue
	client     *core.PoleVpnClient
	networkmgr core.NetworkManager
	device     *core.TunDevice
}

func NewRequestHandler() *RequestHandler {

	return &RequestHandler{mutex: &sync.Mutex{}, status: CLIENT_STOPPED}
}

func (rh *RequestHandler) OnClientEvent(event int, client *core.PoleVpnClient, av *anyvalue.AnyValue) {

	defer core.PanicHandler()

	switch event {
	case core.CLIENT_EVENT_ADDRESS_ALLOCED:
		{
			rh.mutex.Lock()
			defer rh.mutex.Unlock()

			rh.status = CLIENT_STARTED

			var err error
			var routes []string

			if rh.server.Get("UseRemoteRouteRules").AsBool() {
				routes = append(routes, av.Get("route").AsStrArr()...)
			}

			if rh.server.Get("LocalRouteRules").AsStr() != "" {
				routes = append(routes, strings.Split(rh.server.Get("LocalRouteRules").AsStr(), "\n")...)
			}

			if rh.server.Get("ProxyDomains").AsStr() != "" {
				ips := GetRouteIpsFromDomain(strings.Split(rh.server.Get("ProxyDomains").AsStr(), "\n"))
				routes = append(routes, ips...)
			}

			glog.Info("route=", routes, ",allocated ip=", av.Get("ip").AsStr(), ",dns=", av.Get("dns").AsStr())

			if runtime.GOOS == "windows" {
				err = rh.device.GetInterface().SetTunNetwork(av.Get("ip").AsStr() + "/30")
				if err != nil {
					glog.Error("set tun network fail,", err)
					client.Stop()
					return
				}
			}
			av.Set("remoteIp", client.GetRemoteIP())
			err = rh.networkmgr.SetNetwork(rh.device.GetInterface().Name(), av.Get("ip").AsStr(), client.GetRemoteIP(), av.Get("dns").AsStr(), routes)
			if err != nil {
				glog.Error("set network fail,", err)
				client.Stop()
				return
			}
			rh.onCallback(anyvalue.New().Set("event", "allocated").Set("data", av.AsMap()))
		}
	case core.CLIENT_EVENT_STOPPED:
		{
			glog.Info("client stoped")
			rh.networkmgr.RestoreNetwork()
			rh.onCallback(anyvalue.New().Set("event", "stoped").Set("data", nil))
			rh.status = CLIENT_STOPPED
		}
	case core.CLIENT_EVENT_RECONNECTED:
		glog.Info("client reconnected")
		rh.onCallback(anyvalue.New().Set("event", "reconnected").Set("data", nil))
	case core.CLIENT_EVENT_RECONNECTING:
		err := rh.networkmgr.RefreshDefaultGateway()
		if err != nil {
			glog.Error("refresh default gateway fail,", err)
		}
		glog.Info("client reconnecting")
		rh.onCallback(anyvalue.New().Set("event", "reconnecting").Set("data", nil))
	case core.CLIENT_EVENT_STARTED:
		glog.Info("client started")

		var err error
		rh.device, err = core.NewTunDevice()
		if err != nil {
			glog.Error("create device fail,", err)
			go client.Stop()
			return
		}

		client.AttachTunDevice(rh.device)
		rh.onCallback(anyvalue.New().Set("event", "started").Set("data", nil))
	case core.CLIENT_EVENT_ERROR:
		glog.Info("Unexception error,", av.Get("error").AsStr())
		rh.onCallback(anyvalue.New().Set("event", "error").Set("data", av.AsMap()))

	default:
		glog.Error("invalid event=", event)
	}

}

func (rh *RequestHandler) onCallback(av *anyvalue.AnyValue) {

	pkt, _ := av.EncodeJson()
	rh.conn.Send(pkt)
}

func (rh *RequestHandler) OnRequest(pkt []byte, conn Conn) {

	defer core.PanicHandler()

	req, err := anyvalue.NewFromJson(pkt)

	if err != nil {
		glog.Error("decode json fail,", err)
		return
	}
	event := req.Get("event").AsStr()

	if event == "start" {
		err := rh.start(req.Get("data"))
		if err != nil {
			rh.onCallback(anyvalue.New().Set("event", "error").Set("data.error", err.Error()))
		}
	} else if event == "stop" {
		rh.stop()
	} else if event == "getlogs" {

		logFilePath := glog.GetLogPath() + string(os.PathSeparator) + GetAppName() + "-" + GetTimeNowDate() + ".log"

		data, err := ioutil.ReadFile(logFilePath)

		if err != nil {
			glog.Error("read log fail,", err)
			return
		}
		rh.onCallback(anyvalue.New().Set("event", "logs").Set("data.logs", string(data)))

	} else if event == "getbytes" {
		var upBytes, downBytes uint64
		if rh.client != nil {
			upBytes, downBytes = rh.client.GetUpDownBytes()
		}
		rh.onCallback(anyvalue.New().Set("event", "bytes").Set("data.UpBytes", upBytes).Set("data.DownBytes", downBytes))
	} else {
		glog.Error("invalid event,", event)
	}

}

func (rh *RequestHandler) OnConnected(conn Conn) {

	if rh.conn != nil {
		rh.conn.Close(true)
	}
	rh.conn = conn
}

func (rh *RequestHandler) OnClosed(conn Conn, proactive bool) {
	glog.Info(conn.String(), " closed")
	if rh.client != nil && !rh.client.IsStoped() {
		rh.client.Stop()
	}
}

func (rh *RequestHandler) start(server *anyvalue.AnyValue) error {
	rh.mutex.Lock()
	defer rh.mutex.Unlock()

	if rh.status != CLIENT_STOPPED {
		return errors.New("client have started")
	}

	rh.status = CLIENT_STARTING

	glog.Info("Connect to ", server.Get("Endpoint").AsStr())

	rh.server = server

	var err error
	rh.client, err = core.NewPoleVpnClient()

	if err != nil {
		return err
	}

	if runtime.GOOS == "darwin" {
		rh.networkmgr = core.NewDarwinNetworkManager()
	} else if runtime.GOOS == "linux" {
		rh.networkmgr = core.NewLinuxNetworkManager()
	} else if runtime.GOOS == "windows" {
		rh.networkmgr = core.NewWindowsNetworkManager()
	} else {
		return errors.New("os platform not support")
	}

	rh.client.SetEventHandler(rh.OnClientEvent)

	go rh.client.Start(server.Get("Endpoint").AsStr(), server.Get("User").AsStr(), server.Get("Password").AsStr(), server.Get("Sni").AsStr(), server.Get("SkipVerifySSL").AsBool())

	return nil
}

func (rh *RequestHandler) stop() error {
	rh.mutex.Lock()
	defer rh.mutex.Unlock()

	if rh.status != CLIENT_STARTED {
		return errors.New("client haven't started")
	}

	if rh.client != nil {
		rh.client.Stop()
	}
	rh.status = CLIENT_STOPPED
	return nil
}
