package linked_list

import (
	"fmt"
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

func TestLinkedList_LPush(t *testing.T) {
	list := New[string, int8]()
	id := 1679600
	for i := 0; i < 20; i++ {
		uuid := strconv.FormatInt(int64(id+i), 10)
		list.LPush(uuid, 0)
	}

	for {
		ent := list.RPeek()
		if ent == nil {
			break
		}

		fmt.Println(ent.Key)
		list.RPop()
	}
}

func TestLinkedList_RPopAt(t *testing.T) {
	list := New[string, int8]()
	id := 1679600
	for i := 0; i < 20; i++ {
		uuid := strconv.FormatInt(int64(id+i), 10)
		list.LPush(uuid, 0)
	}

	fmt.Println(list.RPopAt(1).Key)
	list = New[string, int8]()
	for i := 0; i < 20; i++ {
		uuid := strconv.FormatInt(int64(id+i), 10)
		list.LPush(uuid, 0)
	}
	fmt.Println(list.RPopAt(19))
}
