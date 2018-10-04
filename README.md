# Temp Gopher

Temp Gopher is a thermostat application written in Go. It is written and tested using a Raspberry Pi, but any hardware platform meeting the requirements will probably work.

## Requirements

You will need a computer (e.g., Raspberry Pi) with the following components:

* A network connection
* DS18B120, 1-wire temperature sensors
* GPIO pins
* Relays for powering on/off your equipment hooked up to the GPIO pins

## Installation

You can build on a Raspberry Pi, however, it can take a long time! I recommend building on a separate computer with a bit more processing power.

```
go get gitea.shoup.io/mike/temp-gopher
cd $GOPATH/src/gitea.shoup.io/mike/temp-gopher
GOOS=linux GOARCH=arm GOARM=6 go build -a -ldflags '-w -s -extldflags "-static"'
scp temp-gopher <raspberrypi>:~/somepath
```

## Configuration

Create a `config.yml` file like this:

```
baseurl: http://<pihostname>:8080 # Base URL to find the app at. Usually your Pi's IP address or hostname, unless using a reverse proxy
sensors:
- id: 28-000008083108 # Id of the DS18b120 sensor
  alias: fermenter # An alias for the sensor
  hightemp: 8 # Maximum temperature you want the sensor to read
  lowtemp: 4 # Minimum tempearture you want the sensor to read
  heatgpio: 5 # GPIO pin the heater is hooked into
  heatinvert: false # Probably false. If true, will set pin to High to turn the heater off
  heatminutes: 1 # Number of minutes below the minimum before the heater turns on
  coolgpio: 17 # GPIO pin the cooler is hooked into
  coolinvert: false # Probably false. If true, will set pin to High to turn the cooler off
  coolminutes: 10 # Number of minutes below the minimum before the cooler turns on
  verbose: false # If true, outputs the current status at every read, approx once per second
```

## Running

You can run it directly in the comment line like:

```
./temp-gopher -c config.yml run
```

You can run it in the background using `nohup`:

```
nohup ./temp-gopher -c config.yml run &
```

Or use `systemctl` or some other process supervisor to run it.

## REST API

There is a very simple REST API for viewing the current configuration and status. The application launches and binds to `:8080`.

To view the current status, query `http://<pi>:8080/api/status`:

```
$ curl -s http://localhost:8080/api/status | jq .
{
  "fermenter": {
    "alias": "fermenter",
    "temp": 19.812,
    "cooling": false,
    "heating": false,
    "reading": "2018-10-03T08:43:05.795870992-06:00",
    "changed": "2999-01-01T00:00:00Z"
  }
}
```

To view the current configuration, query `http://<pi>:8080/api/config`:

```
$ curl -s http://localhost:8080/api/config | jq .
{
  "Sensors": [
    {
      "id": "28-000008083108",
      "alias": "fermenter",
      "hightemp": 30,
      "lowtemp": 27,
      "heatgpio": 13,
      "heatinvert": false,
      "heatminutes": 1,
      "coolgpio": 19,
      "coolinvert": false,
      "coolminutes": 4,
      "verbose": false
    }
  ]
}
```
