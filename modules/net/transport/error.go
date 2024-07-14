package transport

import (
	"errors"
	"fmt"
	"strings"
)

/*
   @Author: orbit-w
   @File: error
   @2023 12月 周六 22:27
*/

var (
	ErrCanceled     = errors.New("context canceled")
	ErrDisconnected = errors.New("disconnected")
	ErrMaxOfRetry   = errors.New(`error_max_of_retry`)
)

func IsClosedConnError(err error) bool {
	/*
		`use of closed file or network connection` (Go ver > 1.8, internal/pool.ErrClosing)
		`mux: listener closed` (cmux.ErrListenerClosed)
	*/
	return err != nil && strings.Contains(err.Error(), "closed")
}

func ExceedMaxIncomingPacket(size uint32) error {
	return errors.New(fmt.Sprintf("exceed max incoming packet size: %d", size))
}

func IsCancelError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "context canceled")
}

func ReceiveBufPutErr(err error) error {
	return errors.New(fmt.Sprintf("receiveBuf put failed: %s", err.Error()))
}

func ReadBodyFailed(err error) error {
	return errors.New(fmt.Sprintf("read body failed: %s", err.Error()))
}

func MaxOfRetryErr(err error) error {
	if err == nil {
		return ErrMaxOfRetry
	}
	return errors.New(fmt.Sprintf("retry failed max limit: %s", err.Error()))
}
