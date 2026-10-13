package spentenergy

import (
	"math"
	"testing"
	"time"
)

func TestDistanceAndMeanSpeedBoundaries(t *testing.T) {
	for _, tc := range []struct {
		name            string
		steps           int
		height          float64
		duration        time.Duration
		distance, speed float64
	}{
		{"normal", 1000, 2, 30 * time.Minute, 0.9, 1.8},
		{"negative steps", -100, 2, time.Hour, 0, 0},
		{"zero steps", 0, 2, time.Hour, 0, 0},
		{"zero duration", 1000, 2, 0, 0.9, 0},
		{"negative duration", 1000, 2, -time.Hour, 0.9, 0},
		{"negative height", 1000, -2, time.Hour, 0, 0},
		{"zero height", 1000, 0, time.Hour, 0, 0},
		{"nan height", 1000, math.NaN(), time.Hour, 0, 0},
		{"infinite height", 1000, math.Inf(1), time.Hour, 0, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := Distance(tc.steps, tc.height); math.IsNaN(got) || math.Abs(got-tc.distance) > 1e-9 {
				t.Fatalf("Distance() = %v, want %v", got, tc.distance)
			}
			if got := MeanSpeed(tc.steps, tc.height, tc.duration); math.IsNaN(got) || math.Abs(got-tc.speed) > 1e-9 {
				t.Fatalf("MeanSpeed() = %v, want %v", got, tc.speed)
			}
		})
	}
}

func TestCaloriesRejectInvalidDurationAndNonFiniteData(t *testing.T) {
	for name, calculate := range map[string]func(int, float64, float64, time.Duration) (float64, error){
		"running": RunningSpentCalories,
		"walking": WalkingSpentCalories,
	} {
		t.Run(name, func(t *testing.T) {
			for _, duration := range []time.Duration{0, -time.Hour} {
				if got, err := calculate(1000, 75, 1.75, duration); err == nil || got != 0 {
					t.Fatalf("invalid duration: got (%v, %v), want (0, error)", got, err)
				}
			}
			for _, invalid := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
				for _, data := range [][2]float64{{invalid, 1.75}, {75, invalid}} {
					if got, err := calculate(1000, data[0], data[1], time.Hour); err == nil || got != 0 {
						t.Fatalf("invalid personal data: got (%v, %v), want (0, error)", got, err)
					}
				}
			}
		})
	}
}
