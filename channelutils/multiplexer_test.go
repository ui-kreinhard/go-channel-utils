package channelutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMux(t *testing.T) {
	inChannel := make(chan bool)

	m := NewChannelMultiplexer(inChannel)
	outChannels := []<-chan bool{
		m.Out(),
		m.Out(),
	}

	inChannel <- true

	assert.True(t, <-outChannels[0])
	assert.True(t, <-outChannels[1])
}
