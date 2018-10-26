package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_PingHandler(t *testing.T) {
	r := gin.New()
	r.GET("/ping", PingHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func Test_ConfigHandler(t *testing.T) {
	testConfig := Config{
		Sensors: []Sensor{
			Sensor{
				Alias: "foo",
			},
		},
		Users:      []User{},
		ListenAddr: ":8080",
	}

	r := gin.New()
	r.GET("/config", ConfigHandler(&testConfig))
	r.GET("/config/sensors/*alias", ConfigHandler(&testConfig))

	// Validate GET request /config
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/config", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	jc, _ := json.Marshal(testConfig)
	assert.Equal(t, string(jc), w.Body.String())

	// Validate GET request to /config/sensors
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/config/sensors/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	jc, _ = json.Marshal(testConfig.Sensors)
	assert.Equal(t, string(jc), w.Body.String())

	// Validate GET request /config/sensors/foo
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/config/sensors/foo", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	jc, _ = json.Marshal(testConfig.Sensors[0])
	assert.Equal(t, string(jc), w.Body.String())

	// Validate not ofund
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/config/sensors/DNE", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func Test_UpdateSensorsHandler(t *testing.T) {
	r := gin.New()
	r.POST("/config/sensors", UpdateSensorsHandler)

	// Test bad request
	buf := bytes.NewBufferString("foobar")
	reader := bufio.NewReader(buf)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/config/sensors", reader)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test good request
	testConfig := Config{
		Sensors: []Sensor{
			Sensor{
				Alias: "foo",
			},
		},
		Users:      []User{},
		ListenAddr: ":8080",
	}
	newSensor := []Sensor{Sensor{Alias: "bar"}}

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

	// Test a POST call
	j, _ := json.Marshal(newSensor)
	buf = bytes.NewBufferString(string(j))
	reader = bufio.NewReader(buf)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/config/sensors", reader)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test internal server error
	configFilePath = "/this/does/not/exist"
	j, _ = json.Marshal(newSensor)
	buf = bytes.NewBufferString(string(j))
	reader = bufio.NewReader(buf)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/config/sensors", reader)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func Test_StatusHandler(t *testing.T) {
	states := make(map[string]State)
	states["foo"] = State{Temp: 5}

	r := gin.New()
	r.GET("/status", StatusHandler(&states))
	r.GET("/status/*alias", StatusHandler(&states))

	// Test all states retrieval
	j, _ := (json.Marshal(states))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/status", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(j), w.Body.String())

	// Test specific state
	j, _ = (json.Marshal(states["foo"]))
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/status/foo", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(j), w.Body.String())

	// Test not found
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/status/DNE", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func Test_JSConfigHandler(t *testing.T) {
	testConfig := Config{
		BaseURL:           "http://localhost:8080",
		DisplayFahrenheit: true,
	}
	jsconfig := "var jsconfig={baseurl:\"" + testConfig.BaseURL +
		"\",fahrenheit:" + strconv.FormatBool(testConfig.DisplayFahrenheit) + "};"

	r := gin.New()
	r.GET("/jsconfig.js", JSConfigHandler(&testConfig))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/jsconfig.js", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsconfig, w.Body.String())
}

func Test_VersionHandler(t *testing.T) {
	type version struct {
		Version string `json:"version"`
	}
	j, _ := json.Marshal(version{Version: Version})

	r := gin.New()
	r.GET("/version", VersionHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/version", nil)

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(j), w.Body.String())
}

func Test_AppHandler(t *testing.T) {
	testConfig := Config{BaseURL: "http://localhost:8080"}
	location := testConfig.BaseURL + "/app/"

	r := gin.New()
	r.Any("/", AppHandler(&testConfig))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusPermanentRedirect, w.Code)
	assert.Equal(t, location, w.Header().Get("Location"))
}

func Test_SetupRouter(t *testing.T) {
	testConfig := Config{
		Sensors: []Sensor{
			Sensor{
				Alias: "foo",
			},
		},
		Users:      []User{},
		ListenAddr: ":8080",
		BaseURL:    "http://localhost:8080",
	}

	states := make(map[string]State)
	states["foo"] = State{}

	// Create a temp file
	tmpfile, err := ioutil.TempFile("", "tempgopher")
	assert.Equal(t, nil, err)
	defer os.Remove(tmpfile.Name()) // Remove the tempfile when done
	configFilePath = tmpfile.Name()

	// Setup a router
	r := SetupRouter(&testConfig, &states)
	assert.IsType(t, gin.New(), r)
}

func Test_GetGinAccounts(t *testing.T) {
	testConfig := Config{
		Users: []User{
			User{
				Name:     "mike",
				Password: "12345",
			},
		},
	}

	testUsers := make(gin.Accounts)
	testUsers["mike"] = "12345"

	actualUsers := GetGinAccounts(&testConfig)

	assert.Equal(t, testUsers, actualUsers)
}
