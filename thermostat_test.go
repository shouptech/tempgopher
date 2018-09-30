package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadTemperature(t *testing.T) {
	data, err := ReadTemperature("28-000008083108")

	assert.Equal(t, 0.0, data)
	assert.NotEqual(t, nil, err)
}
