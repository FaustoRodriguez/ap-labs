package main

import (
	"fmt"
	"time"
)

func main() {
	var commsPerSecond int
	commsPerSecond = 0
	ping := make(chan int)
	pong := make(chan int)
	go play(ping, pong)
	go play(pong, ping)

	for i := 0; i < 1; i++ {
		ping <- 0
		time.Sleep(time.Duration(1) * time.Second)
		commsPerSecond += <-ping
	}

	fmt.Println("Communications Per Second : ", commsPerSecond)
}

func play(ping chan int, pong chan int) {
	for {
		pong <- (1 + <-ping)
	}
}
