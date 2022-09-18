package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ui-kreinhard/go-channel-utils/channelutils"
	"github.com/ui-kreinhard/go-channel-utils/sinks"
)

type Methods struct {
}

func (m *Methods) Map(input int) int {
	fmt.Println("input", input)
	return input
}

func (m *Methods) Filter(input int) bool {
	return input <= 30
}

func (m *Methods) Each(input int) {
	fmt.Println("output", input)
}

func main() {
	randomNumbers := make(chan int)
	go func() {
		for {
			r := rand.Intn(120)
			randomNumbers <- r
			time.Sleep(1 * time.Second)
		}
	}()

	mux := channelutils.NewMap[int, int](randomNumbers, &Methods{}).
		Filter(&Methods{}).Mux()

	mux.Each(&Methods{})
	sinks.NewWebsocketSink(mux.Out(), "localhost:8080")

	for {
		select {}
	}
}
