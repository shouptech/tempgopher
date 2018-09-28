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
	HeatPullup  bool    `yaml:"heatpullup"`
	HeatMinutes float64 `yaml:"heatminutes"`
	CoolGPIO    int32   `yaml:"coolgpio"`
	CoolPullup  bool    `yaml:"coolpullup"`
	CoolMinutes float64 `yaml:"coolminutes"`
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
