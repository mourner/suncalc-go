package suncalc

import (
	"testing"
	"time"
)

var now = time.Now()

func BenchmarkSunPosition(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SunPosition(now, 50.5, 30.5)
	}
}

func BenchmarkSunTimes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SunTimes(now, 50.5, 30.5)
	}
}
