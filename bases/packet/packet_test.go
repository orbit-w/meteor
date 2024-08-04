package packet

import (
	"fmt"
	"testing"
)

/*
   @Author: orbit-w
   @File: packet_test
   @2024 8月 周日 20:36
*/

func Test_Pool(t *testing.T) {
	p := NewPool(maxSize)
	bp := p.GetWithSize(1)
	fmt.Println(bp)
}
