package main

import (
	"fmt"

	"github.com/yryz/ds18b20"
)

func main() {
	sensors, err := ds18b20.Sensors()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Sensors: %v\n", sensors)

	for _, sensor := range sensors {
		t, err := ds18b20.Temperature(sensor)
		if err == nil {
			fmt.Printf("Sensor: %s, Temperature %.2f°C\n", sensor, t)
		}
	}
}
