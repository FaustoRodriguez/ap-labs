package main

import (
	"fmt"
	"time"
)

func main() {
	var maxStages int
	nPipelines := 100000
	first := newPipeline()
	last := first
	for i := 0; i < nPipelines; i++ {
		maxStages++
		previous := last
		newLast := newPipeline()
		go transfer(previous, newLast)
		last = newLast
	}
	startTime := time.Now()
	first <- 1
	<-last
	finnishTime := time.Now()
	transitTime := finnishTime.Sub(startTime)
	fmt.Println("Maximum number of pipeline stages   : ", maxStages)
	fmt.Println("Time to transit trough the pipeline : ", transitTime)
}

func newPipeline() chan int {
	return make(chan int)
}

func transfer(previous chan int, newLast chan int) {
	for i := range previous {
		newLast <- i
	}
}
