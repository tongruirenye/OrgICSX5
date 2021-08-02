package ics

import (
	"testing"
	"time"
)

func TestICSTrigger(t *testing.T) {
	lastTime := time.Now().Add(-5 * time.Hour)
	if !triggerTimer(lastTime, 18, 0) {
		t.Error("trigger wrong")
	}

	lastTime1 := time.Now()
	if triggerTimer(lastTime1, 18, 0) {
		t.Error("trigger1 wrong")
	}
}
