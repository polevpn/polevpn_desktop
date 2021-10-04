package main

type ErrorCode struct {
	Code int
	Msg  string
}

type ReqConnectAccessServer struct {
	AccessServer
}

type RespConnectAccessServer struct {
	ErrorCode
}

type ReqStopAccessServer struct {
}

type RespStopAccessServer struct {
	ErrorCode
}

type ReqAddAccessServer struct {
	AccessServer
}

type RespAddAccessServer struct {
	ErrorCode
}

type ReqGetAllAccessServer struct {
}

type RespGetAllAccessServer struct {
	Servers []AccessServer
	ErrorCode
}

type ReqUpdateAccessServer struct {
	AccessServer
}

type RespUpdateAccessServer struct {
	ErrorCode
}

type ReqDeleteAccessServer struct {
	ID uint
}

type RespDeleteAccessServer struct {
	ErrorCode
}

type ReqGetAllLogs struct {
}

type RespGetAllLogs struct {
	Logs []Logs
	ErrorCode
}
