package timer

import (
	"github.com/google/uuid"
	"sync"
	"time"
)

type TimerCallback func(interface{})

type Timer struct {
	doneChan chan bool
	eventMap map[uuid.UUID]*TimerEvent
	lock     sync.Mutex
}

type TimerEvent struct {
	callback TimerCallback
	capture  interface{}
}

func NewTimer() *Timer {
	timer := &Timer{
		doneChan: make(chan bool),
		eventMap: make(map[uuid.UUID]*TimerEvent),
	}

	return timer
}

func (timer *Timer) startTimerForEvent(eventID uuid.UUID, numRepeats int, timeoutMs uint64) {
	for i := 0; i < numRepeats; i++ {
		select {
		case <-timer.doneChan:
			return
		case <-time.After(time.Duration(timeoutMs) * time.Millisecond):
			timer.lock.Lock()
			event, eventExists := timer.eventMap[eventID]
			if !eventExists {
				timer.lock.Unlock()
				return
			}

			event.callback(event.capture)
			timer.lock.Unlock()
		}
	}

	timer.lock.Lock()
	delete(timer.eventMap, eventID)
	timer.lock.Unlock()
}

func (timer *Timer) AddRepeatingEvent(callback TimerCallback, capture interface{}, timeoutMs uint64, numberRepeats int) uuid.UUID {
	id := uuid.New()
	event := &TimerEvent{
		callback: callback,
		capture:  capture,
	}

	timer.lock.Lock()
	defer timer.lock.Unlock()
	timer.eventMap[id] = event

	go timer.startTimerForEvent(id, numberRepeats, timeoutMs)
	return id
}

func (timer *Timer) RemoveEvent(eventID uuid.UUID) {
	timer.lock.Lock()
	defer timer.lock.Unlock()
	delete(timer.eventMap, eventID)
}

func (timer *Timer) NumEvents() int {
	timer.lock.Lock()
	defer timer.lock.Unlock()
	return len(timer.eventMap)
}

func (timer *Timer) Stop() {
	close(timer.doneChan)
}
