package main

import (
	"sync"

	"github.com/alexflint/go-arg"
	"github.com/stianeikeland/go-rpio"
)

// Version is the current code version of tempgopher
const Version = "0.3.0-dev"

func main() {
	var args struct {
		Action     string `arg:"required,positional" help:"run config"`
		ConfigFile string `arg:"-c,required" help:"path to config file"`
	}

	p := arg.MustParse(&args)
	if args.Action != "run" && args.Action != "config" {
		p.Fail("ACTION must be run or config")
	}

	if args.Action == "config" {
		ConfigCLI(args.ConfigFile)
		return
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
	RunWeb(args.ConfigFile, sc, &wg)

	// Wait for all threads to stop
	wg.Wait()

	// Close the GPIO access
	rpio.Close()
}
