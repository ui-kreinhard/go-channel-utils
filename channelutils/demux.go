package channelutils

type Demux[T any] struct {
	inChannels []chan T
	outChannel chan T
}

func NewDemux[T any](inChannels []chan T) *Demux[T] {
	ret := &Demux[T]{
		inChannels,
		make(chan T),
	}
	go ret.listen()
	return ret
}

func (d *Demux[T]) listen() {
	for _, inChannel := range d.inChannels {
		go d.shovel(inChannel)
	}
}

func (d *Demux[T]) shovel(channel <-chan T) {
	for msg := range channel {
		d.outChannel <- msg
	}
}

func (d *Demux[T]) Out() <-chan T {
	return d.outChannel
}
