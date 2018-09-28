package main

import (
	"errors"
	"log"
	"time"

	"github.com/yryz/ds18b20"
)

// State represents the current state of the thermostat
type State struct {
	Temp     float64
	Cooling  bool
	Heating  bool
	Duration time.Duration
}

// ReadTemperature will return the current temperature (in degrees celsius) of a specific sensor.
func ReadTemperature(id string) (float64, error) {
	sensors, err := ds18b20.Sensors()
	if err != nil {
		return 0.0, err
	}

	for _, sensor := range sensors {
		if sensor != id {
			continue
		}
		return ds18b20.Temperature(sensor)
	}

	return 0.0, errors.New("Sensor not found")
}

// RunThermostat monitors the temperature of the supplied sensor and does its best to keep it at the desired state.
func RunThermostat(sensor Sensor) {
	var s State

	for {
		t, err := ReadTemperature(sensor.ID)
		if err != nil {
			log.Panicln(err)
		}

		s.Temp = t
		log.Printf("Temp: %.2f, Cooling: %t, Heating: %t Duration: %d", s.Temp, s.Cooling, s.Heating, s.Duration)
	}
}
