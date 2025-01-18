// Copyright 2025 Stanislav Tutunaru. All rights reserved.

// Package ftracker is programm module for fitness tracker device.
//
// It provides to get an information from sensors and calculate spent calories during training session.
package ftracker

import (
	"fmt"
	"math"
)

// Main constants needs for calculations.
const (
	lenStep   = 0.65  // mean step length.
	mInKm     = 1000  // count of metres in kilometer.
	minInH    = 60    // count of minutes in one hour.
	kmhInMsec = 0.278 // coefficient for converting km/h to m/s.
	cmInM     = 100   // count of centimeters in one meter.
)

// distance returns the distance (in kilometers) covered by the user during the training session.
//
// *** parameters ***
//
// action int — the number of actions performed (number of steps when walking and running, or strokes when swimming).
func distance(action int) float64 {
	return float64(action) * lenStep / mInKm
}

// meanSpeed returns the value of the average speed during the training session.
//
// *** parameters ***
//
// action int — the number of actions performed (number of steps when walking and running, or strokes when swimming).
//
// duration float64 — training duration in hours.
func meanSpeed(action int, duration float64) float64 {
	if duration == 0 {
		return 0
	}
	distance := distance(action)

	return distance / duration
}

// ShowTrainingInfo returns a string with information about the training.
//
// *** parameters ***
//
// action int — the number of actions performed (number of steps when walking and running, or strokes when swimming).
//
// trainingType string — training type (running, walking, swimming).
//
// duration float64 — training duration in hours.
//
// weight float64 — user weight in kg.
//
// height float64 — user height in m.
//
// lengthPool int — pool length in meters.
//
// countPool int — how many times the user swam across the pool.
func ShowTrainingInfo(action int, trainingType string, duration, weight, height float64, lengthPool, countPool int) string {
	switch {
	case trainingType == "Бег":
		distance := distance(action)
		speed := meanSpeed(action, duration)
		calories := RunningSpentCalories(action, weight, duration)
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", trainingType, duration, distance, speed, calories)
	case trainingType == "Ходьба":
		distance := distance(action)
		speed := meanSpeed(action, duration)
		calories := WalkingSpentCalories(action, duration, weight, height)
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", trainingType, duration, distance, speed, calories)
	case trainingType == "Плавание":
		distance := distance(action)
		speed := swimmingMeanSpeed(lengthPool, countPool, duration)
		calories := SwimmingSpentCalories(lengthPool, countPool, duration, weight)
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", trainingType, duration, distance, speed, calories)
	default:
		return "неизвестный тип тренировки"
	}
}

// Constants for calculating calories consumed during running.
const (
	runningCaloriesMeanSpeedMultiplier = 18   // average velocity multiplier.
	runningCaloriesMeanSpeedShift      = 1.79 // the average number of calories burned while running.
)

// RunningSpentCalories returns the number of calories spent while running.
//
// *** parameters ***
//
// action int — the number of actions performed (number of steps when walking and running, or strokes when swimming).
//
// weight float64 — user weight in kg.
//
// duration float64 — training duration in hours.
func RunningSpentCalories(action int, weight, duration float64) float64 {
	return ((runningCaloriesMeanSpeedMultiplier * meanSpeed(action, duration) * runningCaloriesMeanSpeedShift) * weight / mInKm * duration * minInH)
}

// Constants for calculating calories consumed while walking.
const (
	walkingCaloriesWeightMultiplier = 0.035 // body mass multiplier.
	walkingSpeedHeightMultiplier    = 0.029 // growth multiplier.
)

// WalkingSpentCalories returns the number of calories expended while walking.
//
// *** parameters ***
//
// action int — the number of actions performed (number of steps when walking and running, or strokes when swimming).
//
// duration float64 — training duration in hours.
//
// weight float64 — user weight in kg.
//
// height float64 — user height in m.
func WalkingSpentCalories(action int, duration, weight, height float64) float64 {
	return ((walkingCaloriesWeightMultiplier*weight + (math.Pow(meanSpeed(action, duration)*kmhInMsec, 2)/(height/cmInM))*walkingSpeedHeightMultiplier*weight) * duration * minInH)
}

// Constants for calculating calories expended during swimming.
const (
	swimmingCaloriesMeanSpeedShift   = 1.1 // the average number of calories burned while swimming relative to speed.
	swimmingCaloriesWeightMultiplier = 2   // swimming weight multiplier.
)

// swimmingMeanSpeed returns the average swimming speed.
//
// *** parameters ***
//
// lengthPool int — pool length in meters.
//
// countPool int — how many times the user swam across the pool.
//
// duration float64 — training duration in hours.
func swimmingMeanSpeed(lengthPool, countPool int, duration float64) float64 {
	if duration == 0 {
		return 0
	}

	return float64(lengthPool) * float64(countPool) / mInKm / duration
}

// SwimmingSpentCalories returns the number of calories expended during swimming.
//
// *** parameters ***
//
// lengthPool int — pool length in meters.
//
// countPool int — how many times the user swam across the pool.
//
// duration float64 — training duration in hours.
//
// weight float64 — user weight in kg.
func SwimmingSpentCalories(lengthPool, countPool int, duration, weight float64) float64 {
	return (swimmingMeanSpeed(lengthPool, countPool, duration) + swimmingCaloriesMeanSpeedShift) * swimmingCaloriesWeightMultiplier * weight * duration
}
