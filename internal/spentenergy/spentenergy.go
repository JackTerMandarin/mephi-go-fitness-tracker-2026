package spentenergy

import (
	"errors"
	"math"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе.
)

// WalkingSpentCalories возвращает расход энергии при ходьбе в килокалориях.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}
	return calories * walkingCaloriesCoefficient, nil
}

// RunningSpentCalories возвращает расход энергии при беге в килокалориях.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть больше нуля")
	}
	if weight <= 0 || math.IsNaN(weight) || math.IsInf(weight, 0) {
		return 0, errors.New("вес должен быть положительным конечным числом")
	}
	if height <= 0 || math.IsNaN(height) || math.IsInf(height, 0) {
		return 0, errors.New("рост должен быть положительным конечным числом")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть больше нуля")
	}
	meanSpeed := MeanSpeed(steps, height, duration)
	return weight * meanSpeed * duration.Minutes() / minInH, nil
}

// MeanSpeed возвращает среднюю скорость в километрах в час.
func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return Distance(steps, height) / duration.Hours()
}

// Distance возвращает дистанцию в километрах; рост задаётся в метрах.
func Distance(steps int, height float64) float64 {
	if steps <= 0 || height <= 0 || math.IsNaN(height) || math.IsInf(height, 0) {
		return 0
	}
	stepLength := height * stepLengthCoefficient
	return float64(steps) * stepLength / mInKm
}
