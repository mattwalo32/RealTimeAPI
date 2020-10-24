package timer

import (
	"testing"
	"time"
)

const (
	TOLERANCE = float32(0.1)
)

type count struct {
	value *int
}

func addToCount(data interface{}) {
	c := data.(count)
	*c.value = *c.value + 1
}

func TestSingleEvent(t *testing.T) {
	timer := NewTimer()
	val := 0
	count := count{&val}
	
	callbackTimes := []int{100, 500, 1000}
	for _,timeout := range callbackTimes {
		prevVal := val
		timer.AddRepeatingEvent(addToCount, count, timeout, 1)
		<-time.After(time.Duration(float32(timeout) * TOLERANCE) * time.Millisecond)
		if val != prevVal {
			t.Fatalf("Timer fired too early")
		}

		<-time.After(time.Duration(float32(timeout) * 2 * (1 - TOLERANCE)) * time.Millisecond)
		if val != prevVal + 1 {
			t.Fatalf("Timer fired too late")
		}
	}
}