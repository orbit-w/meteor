package blockreceiver

import (
	"errors"
	"fmt"
)

var (
	ErrCanceled = errors.New("context canceled")
)

func ReceiveBufPutErr(err error) error {
	return errors.New(fmt.Sprintf("receiveBuf put failed: %s", err.Error()))
}
