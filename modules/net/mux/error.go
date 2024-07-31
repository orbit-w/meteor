package mux

import (
	"errors"
	"fmt"
)

/*
   @Author: orbit-w
   @File: error
   @2024 4月 周日 23:32
*/

var (
	ErrCancel   = errors.New("transport_err code: context canceled")
	ErrConnDone = errors.New("error_the_conn_is_done")
)

func NewStreamBufSetErr(err error) error {
	return errors.New(fmt.Sprintf("NewStream set failed: %s", err.Error()))
}

func NewDecodeErr(err error) error {
	return errors.New(fmt.Sprintf("decode data failed: %s", err.Error()))
}
