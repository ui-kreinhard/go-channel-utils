package channelutils

type ChannelMultiplexer[T any] struct {
	inChannel   <-chan T
	OutChannels []chan T
}

func NewChannelMultiplexer[T any](inChannel <-chan T) *ChannelMultiplexer[T] {
	ret := &ChannelMultiplexer[T]{
		inChannel,
		make([]chan T, 0),
	}
	go ret.listen()
	return ret
}

func (p *ChannelMultiplexer[T]) listen() {
	for msg := range p.inChannel {
		for _, o := range p.OutChannels {
			o <- msg
		}
	}
	p.closeAllOutChannels()
}

func (p *ChannelMultiplexer[T]) closeAllOutChannels() {
	for _, o := range p.OutChannels {
		close(o)
	}
}

func (p *ChannelMultiplexer[T]) Out() <-chan T {
	ret := make(chan T)
	p.OutChannels = append(p.OutChannels, ret)
	return ret
}

func (p *ChannelMultiplexer[T]) Filter(filterElement FilterElement[T]) *Filter[T] {
	return NewFilter(p.Out(), filterElement)
}

func (p *ChannelMultiplexer[T]) Each(eachElement ForEachElement[T]) *ForEach[T] {
	return NewForEach(p.Out(), eachElement)
}
