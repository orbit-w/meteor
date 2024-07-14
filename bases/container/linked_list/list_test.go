package linked_list

import (
	"strconv"
	"testing"
)

/*
   @Time: 2023/10/14 21:38
   @Author: david
   @File: list_test
*/

func TestPush(t *testing.T) {
	list := New[string, int8]()
	id := 1679600
	for i := 0; i < 10000; i++ {
		uuid := strconv.FormatInt(int64(id+i), 10)
		list.LPush(uuid, 0)
	}

}
