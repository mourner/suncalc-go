package suncalc

import (
	"math"
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

func degrees(v float64) float64 {
	return v * 180 / math.Pi
}

func TestSunOverhead(t *testing.T) {
	tm, _ := time.Parse(
		time.RFC3339,
		"2012-06-22T12:00:00+00:00")
	azRad, elRad := SunPosition(tm, 55, -3)
	az := degrees(azRad) + 180
	el := degrees(elRad)
	if el < 57 || el > 59 || az < 170 || az > 190 {
		t.Errorf("Out of range sun position for midday June %v,%v", az, el)
	}
}
