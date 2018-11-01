package main

import (
	"errors"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v2"
)

// Influx defines an Influx database configuration
type Influx struct {
	Addr               string  `json:"string"             yaml:"addr"`
	Username           string  `json:"username"           yaml:"username"`
	Password           string  `json:"-"                  yaml:"password"`
	UserAgent          string  `json:"useragent"          yaml:"useragent"`
	Timeout            float64 `json:"timeout"            yaml:"timeout"`
	InsecureSkipVerify bool    `json:"insecureskipverify" yaml:"insecureskipverify"`
	Database           string  `json:"database" yaml:"database"`
}

// Sensor defines configuration for a temperature sensor.
type Sensor struct {
	ID          string  `json:"id"          yaml:"id"`
	Alias       string  `json:"alias"       yaml:"alias"`
	HighTemp    float64 `json:"hightemp"    yaml:"hightemp"`
	LowTemp     float64 `json:"lowtemp"     yaml:"lowtemp"`
	HeatDisable bool    `json:"heatdisable" yaml:"heatdisable"`
	HeatGPIO    int32   `json:"heatgpio"    yaml:"heatgpio"`
	HeatInvert  bool    `json:"heatinvert"  yaml:"heatinvert"`
	HeatMinutes float64 `json:"heatminutes" yaml:"heatminutes"`
	CoolDisable bool    `json:"cooldisable" yaml:"cooldisable"`
	CoolGPIO    int32   `json:"coolgpio"    yaml:"coolgpio"`
	CoolInvert  bool    `json:"coolinvert"  yaml:"coolinvert"`
	CoolMinutes float64 `json:"coolminutes" yaml:"coolminutes"`
	Verbose     bool    `json:"verbose"     yaml:"verbose"`
}

// User defines a user's configuration
type User struct {
	Name     string `json:"name" yaml:"name"`
	Password string `json:"password" yaml:"password"`
}

// Config contains the applications configuration
type Config struct {
	Sensors           []Sensor `yaml:"sensors"`
	Users             []User   `yaml:"users"`
	BaseURL           string   `yaml:"baseurl"`
	ListenAddr        string   `yaml:"listenaddr"`
	DisplayFahrenheit bool     `yaml:"displayfahrenheit"`
	Influx            Influx   `yaml:"influx"`
}

var configFilePath string

// UpdateSensorConfig updates the configuration of an individual sensor and writes to disk
func UpdateSensorConfig(s Sensor) error {
	config, err := LoadConfig(configFilePath)
	if err != nil {
		return err
	}

	for i := range config.Sensors {
		if config.Sensors[i].ID == s.ID {
			copier.Copy(&config.Sensors[i], &s)
		}
	}

	if err = SaveConfig(configFilePath, *config); err != nil {
		return err
	}

	if err = SignalReload(); err != nil {
		return err
	}

	return nil
}

// SignalReload sends a SIGHUP to the process, initiating a configuration reload
func SignalReload() error {
	p := os.Process{Pid: os.Getpid()}
	return p.Signal(syscall.SIGHUP)
}

// SaveConfig will write a new configuration file
func SaveConfig(path string, config Config) error {
	d, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(path, d, 0644); err != nil {
		return err
	}

	return nil
}

// LoadConfig will loads a file and parses it into a Config struct
func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	configFilePath = path
	var config Config
	yaml.Unmarshal(data, &config)

	// Set a default listen address if not define.
	if config.ListenAddr == "" {
		config.ListenAddr = ":8080"
	}

	// Check for Duplicates
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
