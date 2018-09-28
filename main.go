package main

import (
	"log"
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

	config, err := LoadConfig(args.ConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	// Prep for GPIO access
	err = rpio.Open()
	if err != nil {
		log.Fatal(err)
	}

	// Launch the thermostat go routines
	var wg sync.WaitGroup
	for _, sensor := range config.Sensors {
		wg.Add(1)
		go RunThermostat(sensor)
	}
	wg.Wait()

	// Close the GPIO access
	rpio.Close()
}
