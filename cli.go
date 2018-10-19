package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/yryz/ds18b20"
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

	// Configure sensors
	sensors, err := ds18b20.Sensors()
	if err != nil {
		fmt.Println("Couldn't find any sensors. Did you enable the 1-wire bus?")
		fmt.Printf("The error was: %s\n", err)
		os.Exit(1)
	}

	for _, sensor := range sensors {
		fmt.Printf("Configure sensor w/ ID: %s\n", sensor)
		fmt.Print("[Y/n]: ")
		choice := ReadInput(reader, "y")
		if strings.ToLower(choice)[0] != 'y' {
			continue
		}

		var s Sensor
		s.ID = sensor

		fmt.Print("Sensor alias: ")
		s.Alias = ReadInput(reader, "")
		if s.Alias == "" {
			panic("Alias cannot be blank")
		}

		fmt.Print("Disable cooling? [false]: ")
		s.CoolDisable, err = strconv.ParseBool(ReadInput(reader, "false"))
		if err != nil {
			panic(err)
		}

		if !s.CoolDisable {

			fmt.Print("High temperature: ")
			s.HighTemp, err = strconv.ParseFloat(ReadInput(reader, ""), 64)
			if err != nil {
				panic(err)
			}

			fmt.Print("Cooling minutes: ")
			s.CoolMinutes, err = strconv.ParseFloat(ReadInput(reader, ""), 64)
			if err != nil {
				panic(err)
			}

			fmt.Print("Cooling GPIO: ")
			resp, err := strconv.ParseInt(ReadInput(reader, ""), 10, 32)
			s.CoolGPIO = int32(resp)
			if err != nil {
				panic(err)
			}

			fmt.Print("Invert cooling switch [false]: ")
			s.CoolInvert, err = strconv.ParseBool(ReadInput(reader, "false"))
			if err != nil {
				panic(err)
			}
		}

		fmt.Print("Disable heating? [false]: ")
		s.HeatDisable, err = strconv.ParseBool(ReadInput(reader, "false"))
		if err != nil {
			panic(err)
		}

		if !s.HeatDisable {

			fmt.Print("Low temperature: ")
			s.LowTemp, err = strconv.ParseFloat(ReadInput(reader, ""), 64)
			if err != nil {
				panic(err)
			}

			fmt.Print("Heating minutes: ")
			s.HeatMinutes, err = strconv.ParseFloat(ReadInput(reader, ""), 64)
			if err != nil {
				panic(err)
			}

			fmt.Print("Heating GPIO: ")
			resp, err := strconv.ParseInt(ReadInput(reader, ""), 10, 32)
			s.HeatGPIO = int32(resp)
			if err != nil {
				panic(err)
			}

			fmt.Print("Invert heating switch [false]: ")
			s.HeatInvert, err = strconv.ParseBool(ReadInput(reader, "false"))
			if err != nil {
				panic(err)
			}
		}

		fmt.Print("Enable verbose logging [false]: ")
		s.Verbose, err = strconv.ParseBool(ReadInput(reader, "false"))
		if err != nil {
			panic(err)
		}

		config.Sensors = append(config.Sensors, s)
	}

	fmt.Println("Write data to an Influx database?")
	fmt.Print("[Y/n]: ")
	choice := ReadInput(reader, "y")
	if strings.ToLower(choice)[0] == 'y' {
		fmt.Print("Influx address [http://influx:8086]: ")
		config.Influx.Addr = ReadInput(reader, "http://influx:8086")

		fmt.Print("Influx Username []: ")
		config.Influx.Username = ReadInput(reader, "")

		fmt.Print("Influx Password []: ")
		config.Influx.Password = ReadInput(reader, "")

		fmt.Print("Influx UserAgent [InfluxDBClient]: ")
		config.Influx.UserAgent = ReadInput(reader, "InfluxDBClient")

		fmt.Print("Influx timeout (in seconds) [30]: ")
		config.Influx.Timeout, err = strconv.ParseFloat(ReadInput(reader, "30"), 64)
		if err != nil {
			panic(err)
		}

		fmt.Print("Influx database []: ")
		config.Influx.Database = ReadInput(reader, "")

		fmt.Print("Enable InsecureSkipVerify? [fasle]: ")
		config.Influx.InsecureSkipVerify, err = strconv.ParseBool(ReadInput(reader, "false"))
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Enable user authentication?")
	fmt.Print("[Y/n]: ")
	choice = ReadInput(reader, "y")
	if strings.ToLower(choice)[0] == 'y' {
		another := true
		for another {
			fmt.Print("Username: ")
			username := ReadInput(reader, "")
			fmt.Print("Password: ")
			password, err := gopass.GetPasswdMasked()
			if err != nil {
				panic(err)
			}
			config.Users = append(config.Users, User{username, string(password)})

			fmt.Print("Add another user? [y/N]: ")
			choice = ReadInput(reader, "n")
			if strings.ToLower(choice)[0] == 'n' {
				another = false
			}
		}
	}

	return config
}

// ConfigCLI prompts the user for configuration and writes to a config file
func ConfigCLI(path string) {
	config := PromptForConfiguration()

	SaveConfig(path, config)
}
