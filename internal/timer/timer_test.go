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

	if len(timer.eventMap) != 0 {
		t.Errorf("Old events were not deleted")
	}
}

func TestRepeatingEvent(t *testing.T) {
	timer := NewTimer()
	timeout := 500
	repetitions := []int{2, 5}
	
	for _, numReps := range repetitions {
		val := 0
		count := count{&val}
		timer.AddRepeatingEvent(addToCount, count, timeout, numReps)
		<-time.After(time.Duration(float32(timeout) * TOLERANCE) * time.Millisecond)

		for rep := 1; rep <= numReps; rep++ {
			<-time.After(time.Duration(timeout) * time.Millisecond)
			if val != rep {
				t.Fatalf("Timer did not fire in time. Expected %v, got %v", rep, val)
			}
		}
	}
	if len(timer.eventMap) != 0 {
		t.Errorf("Old events were not deleted")
	}
}

func TestRemoveEvent(t *testing.T) {
	timer := NewTimer()
	timeout := 200
	val := 0
	count := count{&val}
	id := timer.AddRepeatingEvent(addToCount, count, timeout, 2)
	timer.RemoveEvent(id)
	<-time.After(time.Duration(float32(timeout) * (1 + TOLERANCE)) * time.Millisecond)
	if val != 0 {
		t.Fatalf("Timer event was not removed")
	}

	if len(timer.eventMap) != 0 {
		t.Errorf("Old events were not deleted")
	}
}