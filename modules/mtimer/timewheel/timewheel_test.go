package timewheel

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeWheel_AddTimer(t *testing.T) {
	//[0,59]
	step := 0
	fmt.Println(60 % 60)
	fmt.Println(60 - step - 1)
	d := time.Minute + time.Second*7
	fmt.Println(d.Milliseconds())
}
