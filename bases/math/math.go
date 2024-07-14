package math

import "github.com/orbit-w/meteor/bases/misc/common"

// PowerOf2 向上取最小的2的平方
// [1, 1<<32]
func PowerOf2(x int) int {
	return 1 << GenericFls(x-1)
}

func GenericFls(x int) int {
	var r int = 32
	if x <= 0 {
		return 0
	}

	if x&0xffff0000 == 0 {
		x <<= 16
		r -= 16
	}

	if x&0xff000000 == 0 {
		x <<= 8
		r -= 8
	}

	if x&0xf0000000 == 0 {
		x <<= 4
		r -= 4
	}

	if x&0xc0000000 == 0 {
		x <<= 2
		r -= 2
	}

	if x&0x80000000 == 0 {
		x <<= 1
		r -= 1
	}
	return r
}

func ABS[T common.Integer](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
