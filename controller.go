package main

import (
	"github.com/polevpn/anyvalue"
	"github.com/polevpn/webview"
)

type Controller struct {
	manager *PoleVPNClientManager
	view    webview.WebView
}

func NewController(view webview.WebView) (*Controller, error) {
	controller := &Controller{view: view}

	manager, err := NewPoleVPNClientManager(controller.EventCallback)

	if err != nil {
		return nil, err
	}
	controller.manager = manager
	return controller, nil
}

func (c *Controller) Bind() {

	c.view.Bind("ConnectAccessServer", c.ConnectAccessServer)
	c.view.Bind("StopAccessServer", c.StopAccessServer)
	c.view.Bind("AddAccessServer", c.AddAccessServer)
	c.view.Bind("UpdateAccessServer", c.UpdateAccessServer)
	c.view.Bind("DeleteAccessServer", c.DeleteAccessServer)
	c.view.Bind("GetAllAccessServer", c.GetAllAccessServer)
	c.view.Bind("GetAllLogs", c.GetAllLogs)
	c.view.Bind("GetUpDownBytes", c.GetUpDownBytes)
	c.view.Bind("GetVersion", c.GetVersion)

}

func (c *Controller) EventCallback(av *anyvalue.AnyValue) {

	json, _ := av.EncodeJson()
	c.view.Dispatch(func() { c.view.Eval("window.onCallback(" + string(json) + ")") })
}

func (c *Controller) ConnectAccessServer(req ReqConnectAccessServer) *RespConnectAccessServer {

	resp := &RespConnectAccessServer{ErrorCode: ErrorCode{Code: 0, Msg: "ok"}}
	err := c.manager.Start(req.AccessServer)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
	}
	return resp
}

func (c *Controller) StopAccessServer(req ReqStopAccessServer) *RespStopAccessServer {
	resp := &RespStopAccessServer{ErrorCode: ErrorCode{Code: 0, Msg: "ok"}}
	err := c.manager.Stop()
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
	}
	return resp
}

func (c *Controller) AddAccessServer(req ReqAddAccessServer) *RespAddAccessServer {

	resp := &RespAddAccessServer{ErrorCode: ErrorCode{Code: 0, Msg: "ok"}}
	err := AddAccessServer(req.AccessServer)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
	}
	return resp
}

func (c *Controller) UpdateAccessServer(req ReqUpdateAccessServer) *RespUpdateAccessServer {

	resp := &RespUpdateAccessServer{ErrorCode: ErrorCode{Code: 0, Msg: "ok"}}
	err := UpdateAccessServer(req.AccessServer)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
	}
	return resp
}

func (c *Controller) DeleteAccessServer(req ReqDeleteAccessServer) *RespDeleteAccessServer {
	resp := &RespDeleteAccessServer{ErrorCode: ErrorCode{Code: 0, Msg: "ok"}}
	err := DeleteAccessServer(req.ID)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
	}
	return resp
}

func (c *Controller) GetAllAccessServer(req ReqGetAllAccessServer) *RespGetAllAccessServer {

	resp := &RespGetAllAccessServer{ErrorCode: ErrorCode{Code: 0, Msg: "ok"}}
	servers, err := GetAllAccessServer()
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
	}
	resp.Servers = servers
	return resp
}

func (c *Controller) GetAllLogs(req ReqGetAllLogs) *RespGetAllLogs {

	resp := &RespGetAllLogs{ErrorCode: ErrorCode{Code: 0, Msg: "ok"}}
	c.manager.GetLogs()
	return resp
}

func (c *Controller) GetUpDownBytes(req ReqGetUpDownBytes) *RespGetUpDownBytes {

	resp := &RespGetUpDownBytes{ErrorCode: ErrorCode{Code: 0, Msg: "ok"}}
	c.manager.GetUpDownBytes()
	return resp
}

func (c *Controller) GetVersion(req ReqGetVersion) *RespGetVersion {

	resp := &RespGetVersion{ErrorCode: ErrorCode{Code: 0, Msg: "ok"}}
	resp.Version = VERSION
	return resp
}
