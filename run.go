package main

import (
	"errors"
	"log"
	"time"

	"github.com/stianeikeland/go-rpio"
	"github.com/yryz/ds18b20"
)

// State represents the current state of the thermostat
type State struct {
	Temp    float64
	Cooling bool
	Heating bool
	Changed time.Time
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

// SetPinState is used to turn a pin on or off.
// If pullup is true, the pin will be pulled up when turned on, and pulled down when turned off.
func SetPinState(pin rpio.Pin, on bool, pullup bool) {
	switch {
	case on && pullup:
		pin.PullUp()
	case on && !pullup:
		pin.PullDown()
	case !on && !pullup:
		pin.PullUp()
	default:
		pin.PullDown()
	}
}

// GetPinState will return true if the pin is on.
// If pullup is true, the pin will be on when pulled up.
func GetPinState(pin rpio.Pin, pullup bool) bool {
	switch pin.Read() {
	case rpio.High:
		return pullup
	default:
		return !pullup
	}
}

// RunThermostat monitors the temperature of the supplied sensor and does its best to keep it at the desired state.
func RunThermostat(sensor Sensor) {
	var s State
	s.Changed = time.Now()

	cpin := rpio.Pin(sensor.CoolGPIO)
	cpin.Output()

	hpin := rpio.Pin(sensor.HeatGPIO)
	hpin.Output()

	for {
		t, err := ReadTemperature(sensor.ID)
		if err != nil {
			log.Panicln(err)
		}

		min := time.Since(s.Changed).Minutes()

		switch {
		case t > sensor.HighTemp && t < sensor.HighTemp:
			log.Panic("Invalid state! Temperature is too high AND too low!")
		case t > sensor.HighTemp && GetPinState(hpin, sensor.HeatPullup):
			SetPinState(hpin, false, sensor.HeatPullup)
		case t > sensor.HighTemp && min > sensor.CoolMinutes:
			SetPinState(cpin, true, sensor.CoolPullup)
		case t > sensor.HighTemp:
			break
		case t < sensor.LowTemp && GetPinState(cpin, sensor.CoolPullup):
			SetPinState(cpin, false, sensor.CoolPullup)
		case t < sensor.LowTemp && min > sensor.HeatMinutes:
			SetPinState(hpin, true, sensor.HeatPullup)
		case t < sensor.LowTemp:
			break
		default: // Turn off both switches
			SetPinState(cpin, false, sensor.CoolPullup)
			SetPinState(hpin, false, sensor.HeatPullup)
		}

		s.Temp = t
		s.Cooling = GetPinState(cpin, sensor.CoolPullup)
		s.Heating = GetPinState(hpin, sensor.HeatPullup)
		log.Printf("Temp: %.2f, Cooling: %t, Heating: %t, Changed: %s", s.Temp, s.Cooling, s.Heating, s.Changed)
	}
}
