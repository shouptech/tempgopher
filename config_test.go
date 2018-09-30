package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadConfig(t *testing.T) {
	testSensor := Sensor{
		ID:          "28-000008083108",
		Alias:       "fermenter",
		HighTemp:    8,
		LowTemp:     4,
		HeatGPIO:    5,
		HeatInvert:  true,
		HeatMinutes: 5,
		CoolGPIO:    17,
		CoolInvert:  false,
		CoolMinutes: 10,
		Verbose:     true,
	}

	testConfig := Config{Sensors: []Sensor{testSensor}}

	loadedConfig, err := LoadConfig("test_config.yml")
	assert.Equal(t, nil, err)
	assert.Equal(t, &testConfig, loadedConfig)

	_, err = LoadConfig("thisfiledoesnotexist")
	assert.NotEqual(t, nil, err)
}
