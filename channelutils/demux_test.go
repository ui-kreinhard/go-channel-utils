package channelutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDemux(t *testing.T) {
	inChannels := []chan int{}
	inChannels = append(inChannels, make(chan int))
	inChannels = append(inChannels, make(chan int))
	inChannels = append(inChannels, make(chan int))

	d := NewDemux(inChannels)
	outChannel := d.Out()

	inChannels[0] <- 0
	inChannels[1] <- 1
	inChannels[2] <- 2

	assert.Equal(t, 0, <-outChannel)
	assert.Equal(t, 1, <-outChannel)
	assert.Equal(t, 2, <-outChannel)
}
