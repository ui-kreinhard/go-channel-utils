package channelutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MapBoolToString struct {
}

func (m *MapBoolToString) Map(in bool) string {
	if in {
		return "true"
	}
	return "false"
}

func TestSimpleMapping(t *testing.T) {
	inChannel := make(chan bool)

	m := NewMap[bool, string](inChannel, &MapBoolToString{})
	outChannel := m.Out()

	inChannel <- true
	assert.Equal(t, "true", <-outChannel)
	inChannel <- false
	assert.Equal(t, "false", <-outChannel)
}
