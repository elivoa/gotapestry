package utils

import (
	"fmt"
	"time"
)

type Timer struct {
	base int // all start time
}

func (t *Timer) Base() int {
	return t.base
}

func (t *Timer) Now() int {
	return time.Now().Nanosecond() - t.base
}

func NewTimer() *Timer {
	return &Timer{base: time.Now().Nanosecond()}
}

func (t *Timer) Log(msg string) {
	fmt.Printf("[TIMER] %v | Use %v ms.\n", msg, (t.Now() / 1000000))
}
