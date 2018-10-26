package main

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadInput(t *testing.T) {
	// Test basic input
	buf := bytes.NewBufferString("foobar\n")
	reader := bufio.NewReader(buf)
	resp := ReadInput(reader, "")
	assert.Equal(t, "foobar", resp)

	// Test default input
	buf = bytes.NewBufferString("\n")
	reader = bufio.NewReader(buf)
	resp = ReadInput(reader, "default")
	assert.Equal(t, "default", resp)

	// Test that a panic occurred
	buf = bytes.NewBufferString("")
	reader = bufio.NewReader(buf)
	assert.Panics(t, func() { ReadInput(reader, "") })
}

func Test_ParsePort(t *testing.T) {
	// Test that 600 is parsed succesfully
	port, err := ParsePort(":600")
	assert.Equal(t, uint16(600), port)
	assert.Equal(t, nil, err)

	// Test failure
	port, err = ParsePort("")
	assert.NotEqual(t, nil, err)
}
