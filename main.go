package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	// run is tracking whether or not the thermostats should run
	run := true

	// done is used to signal the web frontend to stop
	done := make(chan bool)

	// Catch SIGTERM and SIGINT
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT)

	go func() {
		<-sig
		run = false
		done <- true
	}()

	sc := make(chan State)

	// Launch the thermostat go routines
	var wg sync.WaitGroup
	for _, sensor := range config.Sensors {
		wg.Add(1)
		go RunThermostat(sensor, sc, &run, &wg)
	}

	// Launch the web frontend
	wg.Add(1)
	RunWeb(sc, done, &wg)

	// Wait for all threads to stop
	wg.Wait()

	// Close the GPIO access
	rpio.Close()
}
