package number_utils

import (
	"github.com/orbit-w/meteor/bases/misc/common"
	"math/rand"
	"time"
)

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

func Max[T common.Integer](a, b T) T {
	if a > b {
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

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt returns a random integer in the range [min, max)
// It is safe for concurrent use by multiple goroutines.
func RandomInt(min, max int) int {
	if min >= max {
		panic("min should be less than max")
	}
	return rand.Intn(max-min) + min
}

// RandomIntS returns a random integer in the range [min, max)
// It is safe for concurrent use by multiple goroutines.
// New rand source is created for each call.
func RandomIntS(min, max int) int {
	if min >= max {
		panic("min should be less than max")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min) + min
}
