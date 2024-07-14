package number_utils

import "github.com/orbit-w/meteor/bases/misc/common"

/*
   @Time: 2023/8/22 00:17
   @Author: david
   @File: utils
*/

func Min[T common.Integer](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func ABS[T common.Integer](v T) T {
	if v < 0 {
		return -v
	} else {
		return v
	}
}
