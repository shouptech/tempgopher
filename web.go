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

	"github.com/gin-gonic/gin"
)

// PingHandler responds to get requests with the message "pong".
func PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// StatusHandler responds to get requests with the current status of a sensor
func StatusHandler(states *map[string]State) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		if c.Param("alias") == "/" {
			c.JSON(http.StatusOK, states)
		} else if val, ok := (*states)[c.Param("alias")[1:]]; ok {
			c.JSON(http.StatusOK, val)
		} else {
			c.JSON(http.StatusNotFound, "Not found")
		}
	}

	return gin.HandlerFunc(fn)
}

// SetupRouter initializes the gin router.
func SetupRouter(states *map[string]State) *gin.Engine {
	r := gin.Default()

	gin.SetMode(gin.ReleaseMode)

	// Ping
	r.GET("/ping", PingHandler)

	// Status
	r.GET("/api/status/*alias", StatusHandler(states))

	return r
}

// RunWeb launches a web server. sc is used to update the states from the Thermostats.
func RunWeb(sc <-chan State, wg *sync.WaitGroup) {
	// Update sensor states when a new state comes back from the thermostat.
	states := make(map[string]State)
	go func() {
		for {
			s := <-sc
			states[s.Alias] = s
		}
	}()

	// Launch the web server
	r := SetupRouter(&states)
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
