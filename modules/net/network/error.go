package network

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
	ErrCanceled = errors.New("context canceled")
)

// IsClosedConnError 判断是否为关闭连接错误
func IsClosedConnError(err error) bool {
	/*
		`use of closed file or network connection` (Go ver > 1.8, internal/pool.ErrClosing)
		`mux: listener closed` (cmux.ErrListenerClosed)
	*/
	return err != nil && strings.Contains(err.Error(), "closed")
}

// ExceedMaxIncomingPacket 超过最大入站数据包大小
func ExceedMaxIncomingPacket(size uint32) error {
	return errors.New(fmt.Sprintf("exceed max incoming packet size: %d", size))
}

// ReadBodyFailed 读取消息体失败
func ReadBodyFailed(err error) error {
	return errors.New(fmt.Sprintf("Recv body failed: %s", err.Error()))
}

func ReceiveBufPutErr(err error) error {
	return errors.New(fmt.Sprintf("receiveBuf put failed: %s", err.Error()))
}
