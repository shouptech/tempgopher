package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// PingHandler responds to get requests with the message "pong".
func PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// ConfigHandler responds to get requests with the current configuration.
func ConfigHandler(config *Config) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		if c.Param("alias") != "/" && c.Param("alias") != "" {
			alias := c.Param("alias")[1:]
			found := false
			for _, v := range config.Sensors {
				if v.ID == alias {
					c.JSON(http.StatusOK, v)
					found = true
				}
			}
			if !found {
				c.String(http.StatusNotFound, "Not found")
			}
		} else {
			c.JSON(http.StatusOK, *config)
		}
	}
	return gin.HandlerFunc(fn)
}

// StatusHandler responds to get requests with the current status of a sensor
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

	// Status
	r.GET("/api/status", StatusHandler(states))
	r.GET("/api/status/*alias", StatusHandler(states))

	// Config
	r.GET("/api/config", ConfigHandler(config))
	r.GET("/api/config/*alias", ConfigHandler(config))

	return r
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
			config, err = LoadConfig(configpath)
			if err != nil {
				log.Panicln(err)
			}
		}
	}()

	// Launch the web server
	r := SetupRouter(config, &states)
	srv := &http.Server{
		Addr:    ":8080",
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
