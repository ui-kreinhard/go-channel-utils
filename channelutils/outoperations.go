package channelutils

type OutOperations[T any] struct {
	outChannel chan T
}

func NewOutOperations[T any]() *OutOperations[T] {
	return &OutOperations[T]{
		make(chan T),
	}
}

func (o *OutOperations[T]) Out() <-chan T {
	return o.outChannel
}

func (o *OutOperations[T]) Filter(filerElement FilterElement[T]) *Filter[T] {
	return NewFilter(o.Out(), filerElement)
}

func (o *OutOperations[T]) Mux() *ChannelMultiplexer[T] {
	return NewChannelMultiplexer(o.Out())
}

func (o *OutOperations[T]) ForEach(forEach ForEachElement[T]) *ForEach[T] {
	return NewForEach[T](o.outChannel, forEach)
}

func (o *OutOperations[T]) Debug(listenAddress string, mode Mode) *DebugWS[T] {
	return NewDebugWS[T](o.outChannel, mode, listenAddress)
}
