package timewheel

import "errors"

/*
   @Author: orbit-w
   @File: errors
   @2024 9月 周一 22:53
*/

var (
	errTickMs = errors.New("tick must be greater than or equal to 1ms")
)
