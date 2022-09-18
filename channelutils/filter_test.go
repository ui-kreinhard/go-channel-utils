package channelutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestFilterTrue struct{}

func (t *TestFilterTrue) Filter(element bool) bool {
	return element
}

func TestPositiveFilter(t *testing.T) {
	inChannel := make(chan bool)
	f := NewFilter[bool](inChannel, &TestFilterTrue{})
	go f.listen()
	outChannel := f.outChannel

	inChannel <- true

	assert.True(t, <-outChannel)
}

func TestNegativeFilter(t *testing.T) {
	inChannel := make(chan bool)
	f := NewFilter[bool](inChannel, &TestFilterTrue{})
	outChannel := f.outChannel
	close(outChannel)

	inChannel <- false

}
