package gerror

import (
	"fmt"
	"strings"
	"sync"
)

/*
   @Author: orbit-w
   @File: error
   @2024 4月 周三 22:40
*/

type Error struct {
	head string
	text string
}

var builderPool = sync.Pool{New: func() any {
	return &strings.Builder{}
}}

func (e *Error) Error() string {
	w := builderPool.Get().(*strings.Builder)
	defer func() {
		w.Reset()
		builderPool.Put(w)
	}()
	w.WriteString("[")
	w.WriteString(e.head)
	w.WriteString("]: ")
	w.WriteString(e.text)
	return w.String()
}

func New(head string, text string) error {
	return &Error{head: head, text: text}
}

func NewF(head string, format string, args ...interface{}) error {
	return &Error{head: head, text: fmt.Sprintf(format, args...)}
}
