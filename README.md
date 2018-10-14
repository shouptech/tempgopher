# TempGopher

TempGopher is a thermostat application written in Go. It turns a Raspberry Pi into a thermostat you can control with your web browser.

## Requirements

You will need a Raspberry Pi with the following components:

* A network connection
* DS18B120, 1-wire temperature sensors
* GPIO pins
* Relays for powering on/off your equipment hooked up to the GPIO pins

## Install

You will need to setup your Raspberry Pi for this to work.

1. Setup your DS18B20 temperature sensor. [Adafruit has a good tutorial for getting them working](https://learn.adafruit.com/adafruits-raspberry-pi-lesson-11-ds18b20-temperature-sensing/hardware)
2. Connect your relay switches to the GPIO pins
3. [Download the install.sh script to your Raspberry Pi](https://gitlab.com/shouptech/tempgopher/-/jobs/artifacts/master/raw/install.sh?job=build)
4. Run the script! The script will download the latest binary and configure the thermostat with some initial values.
5. After configuration, point your web browser to the configured URL.

## Configuration Script

You will be asked some questions during the initial configuration of TempGopher. You also see some defaults in brackets. If the brackets look good, just hit enter.

* `Listen address?` - The address & port TempGopher should listen on. Omitting the address and just specifying the port means it listens on all addresses. Default is `:8080`.
* `Base URL?` - This is what you will type in to your web browser to access TempGopher. If you don't have DNS configured, should probably be `http://<ipaddress>:8080`
* `Display temperature in fahrenheit?` - Set to true if you want fahrenheit, otherwise defaults to celsius.
* `Configure sensor w/ ID: 28-xxxxx` - If you set up your DS18B20 sensors correctly, you should see it's ID listed. Enter `Y` and answer the prompts to configure it. If you have multiple sensors, you will be asked this question multiple times.
* `Sensor alias:` - Name to display in the web browser for this sensor.
* `High temperature:` - The high temperature to kick the cooling on.
* `Cooling minutes:` - The number of minutes to run the cooler once the temperature is below the High temperature threshold.
* `Cooling GPIO:` - The pin your cooling relay switch is hooked into.
* `Invert cooling switch` - If set to `true`, the cooling will be ON when the switch is LOW. This should usually be `false`, so that is the default
* `Low temperature:` - The low temperature to kick the heating on.
* `Heating minutes:` - The number of minutes to run the heater once the temperature is below the Low temperature threshold.
* `Heating GPIO:` - The pin your heating relay switch is hooked into.
* `Invert heating switch` - If set to `true`, the heating will be ON when the switch is LOW. This should usually be `false`, so that is the default
* `Enable verbose logging` - If set to `true`, TempGopher will display in the console every thermostat reading. This can be quite verbose, so the default is `false`.
* `Write data to an Influx database?` - Whether or not to configure an Influx database
* `Enable user authentication?` - Whether or not to enable authentication

## Example configuration script

```
$ bash install.sh

  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   147  100   147    0     0    357      0 --:--:-- --:--:-- --:--:--   356

TempGopher v0.2.0
You will now be asked a series of questions to help configure your thermostat.
Don't worry, it will all be over quickly.

Default values will be in brackets. Just press enter if they look good.
=====
Listen address?
[:8080]:
Base URL? (This is what you type into your browser to get to the web UI)
[http://beerpi:8080]: http://10.30.14.130:8080
Display temperatures in fahrenheit? (Otherwise uses celsius)
[true]:
Configure sensor w/ ID: 28-000008083108
[Y/n]:
Sensor alias: fermenter
High temperature: 20
Cooling minutes: 4
Cooling GPIO: 19
Invert cooling switch [false]:
Low temperature: 19
Heating minutes: 0.5
Heating GPIO: 13
Invert heating switch [false]:
Enable verbose logging [false]:
Write data to an Influx database?
[Y/n]: y
Influx address [http://influx:8086]:
Influx Username []:
Influx Password []:
Influx UserAgent [InfluxDBClient]:
Influx timeout (in seconds) [30]:
Influx database []: tempgopher
Enable InsecureSkipVerify? [fasle]:
Username: mike
Password: ********
Add another user? [y/N]: y
Username: foo
Password: ***
Add another user? [y/N]: n

```
