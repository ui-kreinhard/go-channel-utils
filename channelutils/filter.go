package channelutils

type Filter[T any] struct {
	inChannel     <-chan T
	outChannel    chan T
	filterElement FilterElement[T]
}

func NewFilter[T any](inChannel <-chan T, filterElement FilterElement[T]) *Filter[T] {
	ret := &Filter[T]{
		inChannel,
		make(chan T),
		filterElement,
	}
	go ret.listen()
	return ret
}

type FilterElement[T any] interface {
	Filter(element T) bool
}

func (f *Filter[T]) listen() {
	for msg := range f.inChannel {
		if f.filterElement.Filter(msg) {
			f.outChannel <- msg
		}
	}
}

func (f *Filter[T]) Out() <-chan T {
	return f.outChannel
}

func (f *Filter[T]) Filter(filerElement FilterElement[T]) *Filter[T] {
	return NewFilter(f.Out(), filerElement)
}

func (f *Filter[T]) Mux() *ChannelMultiplexer[T] {
	return NewChannelMultiplexer(f.Out())
}
