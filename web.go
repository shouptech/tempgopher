package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
	"github.com/jinzhu/copier"
)

// PingHandler responds to GET requests with the message "pong".
func PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// ConfigHandler responds to GET requests with the current configuration.
func ConfigHandler(config *Config) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		if c.Param("alias") != "/" && c.Param("alias") != "" {
			alias := c.Param("alias")[1:]
			found := false
			for _, v := range config.Sensors {
				if v.Alias == alias {
					c.JSON(http.StatusOK, v)
					found = true
				}
			}
			if !found {
				c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			}
		} else if c.Param("alias") == "/" {
			c.JSON(http.StatusOK, config.Sensors)
		} else {
			config.Users = nil // Never return the users in GET requests
			c.JSON(http.StatusOK, config)
		}
	}
	return gin.HandlerFunc(fn)
}

// UpdateSensorsHandler responds to POST requests by updating the stored configuration and issuing a reload to the app
func UpdateSensorsHandler(c *gin.Context) {
	var sensors []Sensor

	if err := c.ShouldBindJSON(&sensors); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, s := range sensors {
		if err := UpdateSensorConfig(s); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// StatusHandler responds to GET requests with the current status of a sensor
func StatusHandler(states *map[string]State) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		if c.Param("alias") == "/" || c.Param("alias") == "" {
			c.JSON(http.StatusOK, states)
		} else if val, ok := (*states)[c.Param("alias")[1:]]; ok {
			c.JSON(http.StatusOK, val)
		} else {
			c.String(http.StatusNotFound, "Not found")
		}
	}

	return gin.HandlerFunc(fn)
}

// GetBox returns a packr.Box object representing the static files.
func GetBox() packr.Box {
	return packr.NewBox("./html")
}

// JSConfigHandler responds to GET requests with the current configuration for the JS app
func JSConfigHandler(config *Config) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		jsconfig := "var jsconfig={baseurl:\"" + config.BaseURL + "\",fahrenheit:" + strconv.FormatBool(config.DisplayFahrenheit) + "};"
		c.String(http.StatusOK, jsconfig)
	}

	return gin.HandlerFunc(fn)
}

// VersionHandler responds to GET requests with the current version of tempgopher
func VersionHandler(c *gin.Context) {
	type version struct {
		Version string `json:"version"`
	}
	c.JSON(http.StatusOK, version{Version: Version})
}

// SetupRouter initializes the gin router.
func SetupRouter(config *Config, states *map[string]State) *gin.Engine {
	// If not specified, put gin in release mode
	if _, ok := os.LookupEnv("GIN_MODE"); !ok {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Midleware
	r.Use(gin.Recovery())
	if gin.Mode() != "release" {
		r.Use(cors.Default())
	} else {
		corsconf := cors.DefaultConfig()
		corsconf.AllowOrigins = []string{config.BaseURL}
		r.Use(cors.New(corsconf))
	}

	// Ping
	r.GET("/ping", PingHandler)

	// API Endpoints
	var api *gin.RouterGroup
	if len(config.Users) == 0 {
		api = r.Group("/api")
	} else {
		api = r.Group("/api", gin.BasicAuth(GetGinAccounts(config)))
	}

	api.GET("/status", StatusHandler(states))
	api.GET("/status/*alias", StatusHandler(states))
	api.GET("/version", VersionHandler)
	api.GET("/config", ConfigHandler(config))
	api.GET("/config/sensors/*alias", ConfigHandler(config))
	api.POST("/config/sensors", UpdateSensorsHandler)

	// App
	r.GET("/jsconfig.js", JSConfigHandler(config))
	r.StaticFS("/app", GetBox())

	// Redirect / to /app
	r.Any("/", func(c *gin.Context) {
		c.Redirect(301, config.BaseURL+"/app/")
	})

	return r
}

// reloadWebConfig reloads the current copy of configuration
func reloadWebConfig(c *Config, p string) error {
	nc, err := LoadConfig(p)
	if err != nil {
		return err
	}

	copier.Copy(&c, &nc)

	return nil
}

// GetGinAccounts returns a gin.Accounts struct with values pulled from a Config struct
func GetGinAccounts(config *Config) gin.Accounts {
	var a gin.Accounts
	for _, user := range config.Users {
		a[user.Name] = user.Password
	}
	return a
}

// RunWeb launches a web server. sc is used to update the states from the Thermostats.
func RunWeb(configpath string, sc <-chan State, wg *sync.WaitGroup) {
	// Update sensor states when a new state comes back from the thermostat.
	states := make(map[string]State)
	go func() {
		for {
			s := <-sc
			states[s.Alias] = s
		}
	}()

	config, err := LoadConfig(configpath)
	if err != nil {
		log.Panicln(err)
	}
	hup := make(chan os.Signal)
	signal.Notify(hup, os.Interrupt, syscall.SIGHUP)
	go func() {
		for {
			<-hup
			err = reloadWebConfig(config, configpath)
			if err != nil {
				log.Panicln(err)
			}
		}
	}()

	// Launch the web server
	r := SetupRouter(config, &states)
	srv := &http.Server{
		Addr:    config.ListenAddr,
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for SIGTERM & SIGINT
	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	signal.Notify(done, os.Interrupt, syscall.SIGINT)
	<-done
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
	wg.Done()
}
