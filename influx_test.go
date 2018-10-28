package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WriteStateToInflux(t *testing.T) {
	state := State{Temp: 32}

	// Test failure with empty config
	err := WriteStateToInflux(state, Influx{})
	assert.NotEqual(t, nil, err)

	// Test failure with missing database
	err = WriteStateToInflux(state, Influx{Addr: "http://127.0.0.1:8086"})
	assert.NotEqual(t, nil, err)

	// Test success with writing to database
	config := Influx{Addr: "http://127.0.0.1:8086", Database: "db"}
	err = WriteStateToInflux(state, config)
	assert.Equal(t, nil, err)
}
