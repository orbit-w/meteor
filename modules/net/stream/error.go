package stream

import "errors"

/*
   @Author: orbit-w
   @File: error
   @2024 4月 周日 23:32
*/

var (
	ErrCancel   = errors.New("transport_err code: context canceled")
	ErrConnDone = errors.New("error_the_conn_is_done")
)
