package main

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Sensor defines configuration for a temperature sensor.
type Sensor struct {
	ID          string  `json:"id"          yaml:"id"`
	Alias       string  `json:"alias"       yaml:"alias"`
	HighTemp    float64 `json:"hightemp"    yaml:"hightemp"`
	LowTemp     float64 `json:"lowtemp"     yaml:"lowtemp"`
	HeatGPIO    int32   `json:"heatgpio"    yaml:"heatgpio"`
	HeatInvert  bool    `json:"heatinvert"  yaml:"heatinvert"`
	HeatMinutes float64 `json:"heatminutes" yaml:"heatminutes"`
	CoolGPIO    int32   `json:"coolgpio"    yaml:"coolgpio"`
	CoolInvert  bool    `json:"coolinvert"  yaml:"coolinvert"`
	CoolMinutes float64 `json:"coolminutes" yaml:"coolminutes"`
	Verbose     bool    `json:"verbose"     yaml:"verbose"`
}

// Config contains the applications configuration
type Config struct {
	Sensors []Sensor `yaml:"sensors"`
}

// LoadConfig will loads a file and parses it into a Config struct
func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	yaml.Unmarshal(data, &config)

	ids := make(map[string]bool)
	aliases := make(map[string]bool)
	for _, v := range config.Sensors {
		if !ids[v.ID] {
			ids[v.ID] = true
		} else {
			return nil, errors.New("Duplicate sensor ID found in configuration")
		}

		if !aliases[v.Alias] {
			aliases[v.Alias] = true
		} else {
			return nil, errors.New("Duplicate sensor alias found in configuration")
		}
	}

	return &config, nil
}
