package main

import (
	"sync"
)

// RunWeb launches a web server
func RunWeb(sc <-chan State, done <-chan bool, wg *sync.WaitGroup) {
	states := make(map[string]State)

	go func() {
		for {
			s := <-sc
			states[s.ID] = s
		}
	}()

	<-done

	wg.Done()
}
