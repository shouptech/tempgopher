package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UpdateSensortConfig(t *testing.T) {
	testConfig := Config{
		Sensors: []Sensor{
			Sensor{
				Alias: "foo",
			},
		},
		Users:      []User{},
		ListenAddr: ":8080",
	}
	newSensor := Sensor{Alias: "bar"}

	// Create a temp file
	tmpfile, err := ioutil.TempFile("", "tempgopher")
	assert.Equal(t, nil, err)

	defer os.Remove(tmpfile.Name()) // Remove the tempfile when done

	configFilePath = tmpfile.Name()

	// Save to tempfile
	err = SaveConfig(tmpfile.Name(), testConfig)
	assert.Equal(t, nil, err)

	// Create a channel to capture SIGHUP
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP)

	// Update the stored config
	UpdateSensorConfig(newSensor)

	// Load the config
	config, err := LoadConfig(tmpfile.Name())
	assert.Equal(t, nil, err)
	assert.Equal(t, "bar", config.Sensors[0].Alias)

	// Validate SIGHUP
	ret := <-sig
	assert.Equal(t, syscall.SIGHUP, ret)
}

func Test_SignalReload(t *testing.T) {
	// Create a channel to capture SIGHUP
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP)

	// Generate SIGHUP
	err := SignalReload()
	assert.Equal(t, nil, err)

	// Validate SIGHUP
	ret := <-sig
	assert.Equal(t, syscall.SIGHUP, ret)
}

func Test_SaveConfig(t *testing.T) {
	// Save zero-valued config
	testConfig := Config{
		Sensors:    []Sensor{},
		Users:      []User{},
		ListenAddr: ":8080",
	}

	// Test writing to a path that doesn't exist
	err := SaveConfig("/this/does/not/exist", testConfig)
	assert.NotEqual(t, nil, err)

	// Create a temp file
	tmpfile, err := ioutil.TempFile("", "tempgopher")
	assert.Equal(t, nil, err)

	defer os.Remove(tmpfile.Name()) // Remove the tempfile when done

	// Save to tempfile
	err = SaveConfig(tmpfile.Name(), testConfig)
	assert.Equal(t, nil, err)

	// Load the config
	config, err := LoadConfig(tmpfile.Name())
	assert.Equal(t, nil, err)
	assert.Equal(t, testConfig, *config)
}

func Test_LoadConfig(t *testing.T) {
	testConfig := Config{
		Sensors: []Sensor{
			Sensor{
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
			},
		},
		Users: []User{
			User{
				Name:     "foo",
				Password: "bar",
			},
		},
		BaseURL:           "https://foo.bar",
		ListenAddr:        ":8080",
		DisplayFahrenheit: true,
		Influx:            Influx{Addr: "http://foo:8086"},
	}

	// Test loading of config
	loadedConfig, err := LoadConfig("tests/test_config.yml")
	assert.Equal(t, nil, err)
	assert.Equal(t, &testConfig, loadedConfig)

	// Test for failures with duplicate IDs and Aliases
	_, err = LoadConfig("tests/duplicate_id.yml")
	assert.NotEqual(t, nil, err)
	_, err = LoadConfig("tests/duplicate_alias.yml")
	assert.NotEqual(t, nil, err)

	// Test for non-existence
	_, err = LoadConfig("DNE")
	assert.NotEqual(t, nil, err)
}
