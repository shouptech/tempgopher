package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	When    time.Time `json:"reading"`
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

// ProcessSensor uses the current temperature and last state to determine if changes need to be made to switches.
func ProcessSensor(sensor Sensor, state State) (State, error) {
	// Read the current temperature
	temp, err := ReadTemperature(sensor.ID)
	if err != nil {
		log.Panicln(err)
	}

	// Initialize the pins
	cpin := rpio.Pin(sensor.CoolGPIO)
	cpin.Output()
	hpin := rpio.Pin(sensor.HeatGPIO)
	hpin.Output()

	// Calculate duration
	duration := time.Since(state.Changed).Minutes()

	switch {
	case temp > sensor.HighTemp && temp < sensor.HighTemp:
		log.Println("Invalid state! Temperature is too high AND too low!")
	case temp > sensor.HighTemp && state.Heating:
		PinSwitch(hpin, false, sensor.HeatInvert)
		state.Heating = false
		state.Changed = time.Now()
	case temp > sensor.HighTemp && state.Cooling:
		break
	case temp > sensor.HighTemp && duration > sensor.CoolMinutes:
		PinSwitch(cpin, true, sensor.CoolInvert)
		state.Cooling = true
		state.Changed = time.Now()
	case temp < sensor.LowTemp && state.Cooling:
		PinSwitch(cpin, false, sensor.CoolInvert)
		state.Cooling = false
		state.Changed = time.Now()
	case temp < sensor.LowTemp && state.Heating:
		break
	case temp < sensor.LowTemp && duration > sensor.HeatMinutes:
		PinSwitch(hpin, true, sensor.HeatInvert)
		state.Heating = true
		state.Changed = time.Now()
	default:
		break
	}

	state.Temp = temp
	if sensor.Verbose {
		log.Printf("%s Temp: %.2f, Cooling: %t, Heating: %t, Duration: %.1f", sensor.Alias, state.Temp, state.Cooling, state.Heating, duration)
	}

	return state, nil
}

// TurnOffSensor turns off all switches for an individual sensor
func TurnOffSensor(sensor Sensor) {
	cpin := rpio.Pin(sensor.CoolGPIO)
	cpin.Output()
	PinSwitch(cpin, false, sensor.CoolInvert)

	hpin := rpio.Pin(sensor.HeatGPIO)
	hpin.Output()
	PinSwitch(hpin, false, sensor.HeatInvert)
}

// TurnOffSensors turns off all sensors defined in the config
func TurnOffSensors(config Config) {
	for _, sensor := range config.Sensors {
		TurnOffSensor(sensor)
	}
}

// RunThermostat monitors the temperature of the supplied sensor and does its best to keep it at the desired state.
func RunThermostat(path string, sc chan<- State, wg *sync.WaitGroup) {
	defer wg.Done()

	// Load Config
	config, err := LoadConfig(path)
	if err != nil {
		log.Panicln(err)
	}

	// Prep for GPIO access
	err = rpio.Open()
	if err != nil {
		log.Panicln(err)
	}
	defer rpio.Close()
	defer TurnOffSensors(*config)

	// Track if thermostats should run
	run := true

	// Start with everything off
	TurnOffSensors(*config)

	// Listen for SIGHUP to reload config
	hup := make(chan os.Signal)
	signal.Notify(hup, os.Interrupt, syscall.SIGHUP)
	go func() {
		for {
			<-hup
			log.Println("Reloading configuration")
			config, err = LoadConfig(path)
			if err != nil {
				log.Panicln(err)
			}
		}
	}()

	// Listen for SIGTERM & SIGINT to quit
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT)
	go func() {
		<-sig
		run = false
	}()

	states := make(map[string]State)
	// For each sensor, run through the thermostat logic
	for run {
		for _, v := range config.Sensors {
			// Create an initial state if there's not one already
			if _, ok := states[v.ID]; !ok {
				state := State{
					Alias:   v.Alias,
					When:    time.Now(),
					Changed: time.Now(),
				}
				states[v.ID] = state
			}

			// Process the sensor
			states[v.ID], err = ProcessSensor(v, states[v.ID])
			if err != nil {
				log.Panicln(err)
			}

			// Write the returned state to the channel (don't block if nothing is available to listen)
			select {
			case sc <- states[v.ID]:
				break
			default:
				break
			}
		}
	}

	log.Println("Shutting down thermostat")
}
