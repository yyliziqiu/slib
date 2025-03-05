package stime

import (
	"time"
)

type Timer struct {
	StartTime time.Time
	PauseTime time.Time
}

func NewTimer() Timer {
	return Timer{
		StartTime: time.Now(),
		PauseTime: time.Now(),
	}
}

func (t Timer) Pause() time.Duration {
	d := time.Now().Sub(t.PauseTime)
	t.PauseTime = time.Now()
	return d
}

func (t Timer) Stop() time.Duration {
	return time.Now().Sub(t.StartTime)
}

func (t Timer) Pauses() string {
	return ManualDuration(t.Pause())
}

func (t Timer) Stops() string {
	return ManualDuration(t.Stop())
}
