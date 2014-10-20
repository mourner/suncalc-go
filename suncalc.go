package main

import (
	"fmt"
	"math"
	"time"
)

const rad = math.Pi / 180

// time conversions

const daySec = 60 * 60 * 24
const j1970 = 2440588.0
const j2000 = 2451545.0

func toJulian(t time.Time) float64 {
	return float64(t.Unix()) / daySec - 0.5 + j1970
}
func fromJulian(j float64) time.Time {
	return time.Unix(int64((j + 0.5 - j1970) * daySec), 0)
}
func toDays(t time.Time) float64 {
	return toJulian(t) - j2000
}

// general utilities for celestial body position

const e = rad * 23.4397

func rightAscension(l, b float64) float64 {
	return math.Atan2(math.Sin(l) * math.Cos(e) - math.Tan(b) * math.Sin(e), math.Cos(l))
}
func declination(l, b float64) float64 {
	return math.Asin(math.Sin(b) * math.Cos(e) + math.Cos(b) * math.Sin(e) * math.Sin(l))
}
func azimuth(H, phi, dec float64) float64 {
	return math.Atan2(math.Sin(H), math.Cos(H) * math.Sin(phi) - math.Tan(dec) * math.Cos(phi))
}
func altitude(H, phi, dec float64) float64 {
	return math.Sin(math.Sin(phi) * math.Sin(dec) + math.Cos(phi) * math.Cos(dec) * math.Cos(H))
}
func siderealTime(d, lw float64) float64 {
	return rad * (280.16 + 360.9856235 * d) - lw
}

// general sun calculations

func solarMeanAnomaly(d float64) float64 {
	return rad * (357.5291 + 0.98560028 * d)
}
func eclipticLongitude(m float64) float64 {
	c := rad * (1.9148 * math.Sin(m) + 0.02 * math.Sin(2 * m) + 0.0003 * math.Sin(3 * m)) // equation of center
	p := rad * 102.9372 // perihelion of the Earth
	return m + c + p + math.Pi
}
func sunCoords(d float64) (float64, float64) {
	l := eclipticLongitude(solarMeanAnomaly(d))
	return declination(l, 0), rightAscension(l, 0)
}

// returns sun's azimuth and altitude given time and latitude/longitude

func SunPosition(t time.Time, lat, lng float64) (float64, float64) {
	lw  := rad * -lng
	phi := rad * lat
	d := toDays(t)
	dec, ra := sunCoords(d)
	h := siderealTime(d, lw) - ra

	return azimuth(h, phi, dec), altitude(h, phi, dec)
}

// calculations for sun times

const j0 = 0.0009

func julianCycle(d, lw float64) float64 {
	return math.Floor(d - j0 - lw / (2.0 * math.Pi) + 0.5)
}
func approxTransit(ht, lw, n float64) float64 {
	return j0 + (ht + lw) / (2.0 * math.Pi) + n
}
func solarTransitJ(ds, m, l float64) float64 {
	return j2000 + ds + 0.0053 * math.Sin(m) - 0.0069 * math.Sin(2 * l)
}
func hourAngle(h, phi, d float64) float64 {
	return math.Acos((math.Sin(h) - math.Sin(phi) * math.Sin(d)) / (math.Cos(phi) * math.Cos(d)))
}

// returns set time for the given sun altitude
func getSetJ(h, lw, phi, dec, n, m, l float64) float64 {
	w := hourAngle(h, phi, dec)
	a := approxTransit(w, lw, n)
	return solarTransitJ(a, m, l)
}

// sun times configuration

type SunAngle struct {
	angle    float64
	riseName string
	setName  string
}

var sunAngles []SunAngle = []SunAngle{
	SunAngle{-0.833, "sunrise", "sunset"},
	SunAngle{-0.3, "sunriseEnd", "sunsetStart"},
	SunAngle{-6.0, "dawn", "dusk"},
	SunAngle{-12.0, "nauticalDawn", "nauticalDusk"},
	SunAngle{-18.0, "nightEnd", "night"},
	SunAngle{6.0, "goldenHourEnd", "goldenHour"},
}

// calculates sun times for a given date and latitude/longitude
func SunTimes(t time.Time, lat, lng float64) map[string]time.Time {
	lw := rad * -lng
	phi := rad * lat

	d := toDays(t)
	n := julianCycle(d, lw)
	ds := approxTransit(0, lw, n)

	m := solarMeanAnomaly(ds)
	l := eclipticLongitude(m)
	dec := declination(l, 0)

	jNoon := solarTransitJ(ds, m, l)

	times := map[string]time.Time{
		"solarNoon": fromJulian(jNoon),
		"nadir": fromJulian(jNoon - 0.5),
	}

	for _, sunAngle := range sunAngles {
		jSet := getSetJ(sunAngle.angle * rad, lw, phi, dec, n, M, L)

		times[sunAngle.riseName] = fromJulian(jNoon - (jSet - jNoon))
		times[sunAngle.setName] = fromJulian(jSet)
	}

	return times
}


func main() {
	azimuth, altitude := SunPosition(time.Now(), 50.5, 30.5)
	fmt.Println("position", azimuth, altitude)

	times := SunTimes(time.Now(), 50.5, 30.5)
	fmt.Println("times", times)
}
