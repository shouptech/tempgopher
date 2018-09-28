package main

import (
	"log"
	"sync"

	arg "github.com/alexflint/go-arg"
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

	var wg sync.WaitGroup

	for _, sensor := range config.Sensors {
		wg.Add(1)
		go RunThermostat(sensor)
	}

	wg.Wait()
}
