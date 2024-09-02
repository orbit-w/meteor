package timewheel

import (
	"fmt"
	"testing"
)

func TestTimeWheel_AddTimer(t *testing.T) {
	//[0,59]
	step := 0
	fmt.Println(60 % 60)
	fmt.Println(60 - step - 1)
}
