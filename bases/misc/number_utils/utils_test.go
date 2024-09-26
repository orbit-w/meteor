package number_utils

import "testing"

func TestRandomInt(t *testing.T) {
	min, max := 1, 100
	for i := 0; i < 1000; i++ {
		v := RandomInt(min, max)
		if v < min || v >= max {
			t.Errorf("RandomInt(%d, %d) = %d, out of range", min, max, v)
		}
	}
}
