package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ReadInput reads the next line from a Reader. It will return 'def' if nothing
// was input, or it will return the response.
func ReadInput(r *bufio.Reader, def string) string {
	resp, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}

	if resp == "\n" {
		return def
	}

	return resp[:len(resp)-1]
}

// ParsePort returns the port number from a listen address
func ParsePort(addr string) (uint16, error) {
	parts := strings.Split(addr, ":")

	port, err := strconv.ParseUint(parts[len(parts)-1], 10, 16)

	if err != nil {
		return 0, err
	}

	return uint16(port), nil
}

// PromptForConfiguration walks a user through configuration
func PromptForConfiguration() Config {
	reader := bufio.NewReader(os.Stdin)

	var config Config

	fmt.Printf("TempGopher v%s\n", Version)
	fmt.Println("You will now be asked a series of questions to help configure your thermostat.")
	fmt.Println("Don't worry, it will all be over quickly.")
	fmt.Println("\nDefault values will be in brackets. Just press enter if they look good.")
	fmt.Println("=====")

	fmt.Print("Listen address?\n[:8080]: ")
	config.ListenAddr = ReadInput(reader, ":8080")

	port, err := ParsePort(config.ListenAddr)
	if err != nil {
		panic(err)
	}
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	fmt.Println("Base URL? (This is what you type into your browser to get to the web UI)")
	defURL := fmt.Sprintf("http://%s:%d", hostname, port)
	fmt.Printf("[%s]: ", defURL)
	config.BaseURL = ReadInput(reader, defURL)

	fmt.Println("Display temperatures in fahrenheit? (Otherwise uses celsius)")
	fmt.Print("[true]: ")

	config.DisplayFahrenheit, err = strconv.ParseBool(ReadInput(reader, "true"))
	if err != nil {
		panic(err)
	}

	return config
}

// ConfigCLI prompts the user for configuration and writes to a config file
func ConfigCLI(path string) {
	config := PromptForConfiguration()

	SaveConfig(path, config)
}
