package timewheel

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCallback_Exec(t *testing.T) {
	t.Run("successful execution", func(t *testing.T) {
		result := 0
		cb := newCallback(func(args ...any) {
			val := args[0].(int)
			result = val * 2
		}, 5)

		cb.Exec()
		assert.Equal(t, 10, result, "callback should multiply input by 2")
	})

	t.Run("multiple arguments", func(t *testing.T) {
		var str string
		var num int
		cb := newCallback(func(args ...any) {
			str = args[0].(string)
			num = args[1].(int)
		}, "test", 42)

		cb.Exec()
		assert.Equal(t, "test", str, "string argument should be passed correctly")
		assert.Equal(t, 42, num, "integer argument should be passed correctly")
	})

	t.Run("panic recovery", func(t *testing.T) {
		cb := newCallback(func(args ...any) {
			panic("test panic")
		})

		// Should not panic
		cb.Exec()
	})

	t.Run("nil function", func(t *testing.T) {
		cb := Callback{
			f:    nil,
			args: []interface{}{1, 2, 3},
		}

		// Should not panic, but log error
		cb.Exec()
	})
}
