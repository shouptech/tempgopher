package main

import (
	"errors"
	"log"
	"os"
	"sync"
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
// If invert is false, the pin will be high turned on, and low when turned off.
func SetPinState(pin rpio.Pin, on bool, invert bool) {
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
func RunThermostat(sensor Sensor, sig <-chan os.Signal, wg *sync.WaitGroup) {
	var s State
	s.Changed = time.Now()

	cpin := rpio.Pin(sensor.CoolGPIO)
	cpin.Output()

	hpin := rpio.Pin(sensor.HeatGPIO)
	hpin.Output()

	SetPinState(cpin, false, sensor.CoolInvert)
	SetPinState(hpin, false, sensor.HeatInvert)

	run := true
	go func() {
		<-sig
		run = false
	}()

	for run {
		t, err := ReadTemperature(sensor.ID)
		if err != nil {
			log.Panicln(err)
		}

		min := time.Since(s.Changed).Minutes()

		switch {
		case t > sensor.HighTemp && t < sensor.HighTemp:
			log.Panic("Invalid state! Temperature is too high AND too low!")
		case t > sensor.HighTemp && s.Heating:
			SetPinState(hpin, false, sensor.HeatInvert)
			log.Printf("%s Turned off heat", sensor.Alias)
			s.Heating = false
			s.Changed = time.Now()
		case t > sensor.HighTemp && s.Cooling:
			break
		case t > sensor.HighTemp && min > sensor.CoolMinutes:
			SetPinState(cpin, true, sensor.CoolInvert)
			log.Printf("%s Turned on cool", sensor.Alias)
			s.Cooling = true
			s.Changed = time.Now()
		case t < sensor.LowTemp && s.Cooling:
			SetPinState(cpin, false, sensor.CoolInvert)
			log.Printf("%s Turned off cool", sensor.Alias)
			s.Cooling = false
			s.Changed = time.Now()
		case t < sensor.LowTemp && s.Heating:
			break
		case t < sensor.LowTemp && min > sensor.HeatMinutes:
			SetPinState(hpin, true, sensor.HeatInvert)
			log.Printf("%s Turned on heat", sensor.Alias)
			s.Heating = true
			s.Changed = time.Now()
		default:
			break
		}

		s.Temp = t
		log.Printf("%s Temp: %.2f, Cooling: %t, Heating: %t, Duration: %.1f", sensor.Alias, s.Temp, s.Cooling, s.Heating, min)
	}

	SetPinState(cpin, false, sensor.CoolInvert)
	SetPinState(hpin, false, sensor.HeatInvert)
	wg.Done()
}
