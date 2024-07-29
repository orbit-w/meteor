package agent_stream

import "errors"

/*
   @Author: orbit-w
   @File: error
   @2024 4月 周日 23:32
*/

var (
	ErrRpcDisconnected  = errors.New("error_rpc_disconnected")
	ErrRpcDisconnectedP = "error_rpc_disconnected"
	ErrConnDone         = errors.New("error_the_conn_is_done")
	ErrMaxOfRetry       = errors.New(`error_max_of_retry`)
	ErrCancel           = errors.New("transport_err code: context canceled")
	ErrStreamShutdown   = errors.New("transport_err code: mux shutdown")
	ErrStreamQuotaEmpty = errors.New("err_stream_quota_empty")
)
