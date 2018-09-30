package main

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/stianeikeland/go-rpio"
	"github.com/yryz/ds18b20"
)

// State represents the current state of the thermostat
type State struct {
	Alias   string    `json:"alias"`
	Temp    float64   `json:"temp"`
	Cooling bool      `json:"cooling"`
	Heating bool      `json:"heating"`
	Changed time.Time `json:"changed"`
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

// PinSwitch is used to turn a pin on or off.
// If invert is false, the pin will be high turned on, and low when turned off.
func PinSwitch(pin rpio.Pin, on bool, invert bool) {
	switch {
	case on && !invert:
		pin.High()
	case on && invert:
		pin.Low()
	case !on && invert:
		pin.High()
	default:
		pin.Low()
	}
}

// RunThermostat monitors the temperature of the supplied sensor and does its best to keep it at the desired state.
func RunThermostat(sensor Sensor, sc chan<- State, run *bool, wg *sync.WaitGroup) {
	var s State
	s.Alias = sensor.Alias
	s.Changed = time.Now()

	cpin := rpio.Pin(sensor.CoolGPIO)
	cpin.Output()

	hpin := rpio.Pin(sensor.HeatGPIO)
	hpin.Output()

	PinSwitch(cpin, false, sensor.CoolInvert)
	PinSwitch(hpin, false, sensor.HeatInvert)

	for *run {
		t, err := ReadTemperature(sensor.ID)
		if err != nil {
			log.Panicln(err)
		}

		min := time.Since(s.Changed).Minutes()

		switch {
		case t > sensor.HighTemp && t < sensor.HighTemp:
			log.Println("Invalid state! Temperature is too high AND too low!")
		case t > sensor.HighTemp && s.Heating:
			PinSwitch(hpin, false, sensor.HeatInvert)
			s.Heating = false
			s.Changed = time.Now()
		case t > sensor.HighTemp && s.Cooling:
			break
		case t > sensor.HighTemp && min > sensor.CoolMinutes:
			PinSwitch(cpin, true, sensor.CoolInvert)
			s.Cooling = true
			s.Changed = time.Now()
		case t < sensor.LowTemp && s.Cooling:
			PinSwitch(cpin, false, sensor.CoolInvert)
			s.Cooling = false
			s.Changed = time.Now()
		case t < sensor.LowTemp && s.Heating:
			break
		case t < sensor.LowTemp && min > sensor.HeatMinutes:
			PinSwitch(hpin, true, sensor.HeatInvert)
			s.Heating = true
			s.Changed = time.Now()
		default:
			break
		}

		s.Temp = t
		if sensor.Verbose {
			log.Printf("%s Temp: %.2f, Cooling: %t, Heating: %t, Duration: %.1f", sensor.Alias, s.Temp, s.Cooling, s.Heating, min)
		}

		select {
		case sc <- s:
			break
		default:
			break
		}
	}

	log.Printf("%s Shutting down thermostat", sensor.Alias)
	PinSwitch(cpin, false, sensor.CoolInvert)
	PinSwitch(hpin, false, sensor.HeatInvert)
	wg.Done()
}
