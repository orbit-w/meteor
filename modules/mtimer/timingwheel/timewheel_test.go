package timewheel

import (
	"fmt"
	"sync/atomic"
	"testing"
)

type MyStruct struct {
	Value int64
}

func TestTimeWheel_Atomic(t *testing.T) {
	var ptr atomic.Pointer[MyStruct]

	// Create a new instance of MyStruct
	newStruct := &MyStruct{Value: 42}
	swapped := ptr.CompareAndSwap(nil, newStruct)
	fmt.Println("Swapped:", swapped)

	// Store the new instance atomically
	ptr.Store(newStruct)

	// Load the value atomically
	loadedStruct := ptr.Load()
	fmt.Println("Loaded Value:", loadedStruct.Value)

	// Compare and Swap
	oldStruct := &MyStruct{Value: 42}
	newStruct2 := &MyStruct{Value: 100}
	swapped = ptr.CompareAndSwap(oldStruct, newStruct2)
	fmt.Println("Swapped:", swapped)

	// Load the new value
	loadedStruct = ptr.Load()
	fmt.Println("New Loaded Value:", loadedStruct.Value)
}
