package main

import (
	"sync"

	"github.com/alexflint/go-arg"
	"github.com/stianeikeland/go-rpio"
)

func main() {
	var args struct {
		Action     string `arg:"required,positional" help:"run"`
		ConfigFile string `arg:"-c,required" help:"path to config file"`
	}

	p := arg.MustParse(&args)
	if args.Action != "run" {
		p.Fail("ACTION must be run")
	}

	// Create a channel for receiving of state
	sc := make(chan State)

	// Use to track running routines
	var wg sync.WaitGroup

	// Launch the thermostat go routines
	wg.Add(1)
	go RunThermostat(args.ConfigFile, sc, &wg)

	// Launch the web frontend
	wg.Add(1)
	RunWeb(sc, &wg)

	// Wait for all threads to stop
	wg.Wait()

	// Close the GPIO access
	rpio.Close()
}
