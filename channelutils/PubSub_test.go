package channelutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublishSubscribe(t *testing.T) {
	inChannel := make(chan bool)

	ps := NewPublishSubscribe(inChannel)
	sub1 := ps.Subscribe()
	sub2 := ps.Subscribe()

	inChannel <- true
	assert.True(t, <-sub1.Out())
	assert.True(t, <-sub2.Out())

}

func TestPublishUnsubscribe(t *testing.T) {
	inChannel := make(chan bool)

	ps := NewPublishSubscribe(inChannel)
	sub1 := ps.Subscribe()

	sub1.CancelSubscription()
	assert.True(t, isClosed(sub1.Out()))
}

func TestPublishSubscribeFilter(t *testing.T) {
	inChannel := make(chan bool)

	ps := NewPublishSubscribe(inChannel)
	sub1 := ps.SubscribeWithFilter(&FilterTrue{})
	sub2 := ps.SubscribeWithFilter(&FilterFalse{})

	close(sub1.out)
	inChannel <- false

	assert.False(t, <-sub2.Out())
}

type FilterTrue struct{}

func (f *FilterTrue) Filter(in bool) bool {
	return in
}

type FilterFalse struct{}

func (f *FilterFalse) Filter(in bool) bool {
	return !in
}

func isClosed[T any](ch <-chan T) bool {
	select {
	case <-ch:
		return true
	default:
	}
	return false
}
