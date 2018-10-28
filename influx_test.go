package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WriteStateToInflux(t *testing.T) {
	influxAddr := os.Getenv("INFLUXDB_ADDR")
	state := State{Temp: 32}

	// Test failure with empty config
	err := WriteStateToInflux(state, Influx{})
	assert.NotEqual(t, nil, err)

	// Test failure with missing database
	err = WriteStateToInflux(state, Influx{Addr: influxAddr})
	assert.NotEqual(t, nil, err)

	// Test success with writing to database
	config := Influx{Addr: influxAddr, Database: "db"}
	err = WriteStateToInflux(state, config)
	assert.Equal(t, nil, err)
}
