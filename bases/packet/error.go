package packet

import "errors"

/*
   @Author: orbit-w
   @File: error
   @2023 11月 周日 15:07
*/

var (
	ErrReadByteFailed        = errors.New("read_byte_failed")
	ErrReadBytesHeaderFailed = errors.New("error_read_bytes_header_failed")
	ErrOutOfRange            = errors.New("error_out_of_range")
)
