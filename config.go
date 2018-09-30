package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Sensor defines configuration for a temperature sensor.
type Sensor struct {
	ID          string  `yaml:"id"`
	Alias       string  `yaml:"alias"`
	HighTemp    float64 `yaml:"hightemp"`
	LowTemp     float64 `yaml:"lowtemp"`
	HeatGPIO    int32   `yaml:"heatgpio"`
	HeatInvert  bool    `yaml:"heatinvert"`
	HeatMinutes float64 `yaml:"heatminutes"`
	CoolGPIO    int32   `yaml:"coolgpio"`
	CoolInvert  bool    `yaml:"coolinvert"`
	CoolMinutes float64 `yaml:"coolminutes"`
	Verbose     bool    `yaml:"verbose"`
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

	return &config, nil
}
