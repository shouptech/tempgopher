package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var states map[string]State

// PingHandler responds to get requests with the message "pong".
func PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// StatusHandler responds to get requests with the current status of a sensor
func StatusHandler(c *gin.Context) {
	if c.Param("alias") == "/" {
		c.JSON(http.StatusOK, states)
	} else if val, ok := states[c.Param("alias")[1:]]; ok {
		c.JSON(http.StatusOK, val)
	} else {
		c.JSON(http.StatusNotFound, "Not found")
	}
}

// SetupRouter initializes the gin router.
func SetupRouter() *gin.Engine {
	r := gin.Default()

	gin.SetMode(gin.ReleaseMode)

	// Ping
	r.GET("/ping", PingHandler)

	// Status
	r.GET("/api/status/*alias", StatusHandler)

	return r
}

// RunWeb launches a web server. sc is used to update the states from the Thermostats.
func RunWeb(sc <-chan State, done <-chan bool, wg *sync.WaitGroup) {
	// Update sensor states when a new state comes back from the thermostat.
	states = make(map[string]State)
	go func() {
		for {
			s := <-sc
			states[s.Alias] = s
		}
	}()

	// Launch the web server
	r := SetupRouter()
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

	// Wait for the done signal
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
