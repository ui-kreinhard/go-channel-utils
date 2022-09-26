package channelutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Out struct {
	number int
	ok     bool
}

type JoinIntBool struct {
}

func (j *JoinIntBool) Join(number int, ok bool) Out {
	return Out{
		number,
		ok,
	}
}

func TestJoin(t *testing.T) {
	in1 := make(chan int)
	in2 := make(chan bool)

	joiner := NewJoin[int, bool, Out](in1, in2, &JoinIntBool{})

	in1 <- 1
	in1 <- 2
	in1 <- 3

	in2 <- true

	joined := <-joiner.Out()
	assert.Equal(t, Out{3, true}, joined)
}
