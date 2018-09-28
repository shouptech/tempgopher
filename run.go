package main

import (
	"errors"
	"log"

	"github.com/yryz/ds18b20"
)

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
	for {
		t, err := ReadTemperature(sensor.ID)

		if err != nil {
			log.Panicln(err)
		}

		log.Printf("%.2f", t)
	}
}
