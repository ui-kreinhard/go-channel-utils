package channelutils

type ForEach[T any] struct {
	inChannel      <-chan T
	forEachElement ForEachElement[T]
}

type ForEachElement[T any] interface {
	Each(element T)
}

func NewForEach[T any](inChannel <-chan T, forEachElement ForEachElement[T]) *ForEach[T] {
	ret := &ForEach[T]{
		inChannel,
		forEachElement,
	}
	go ret.listen()
	return ret
}

func (f *ForEach[T]) listen() {
	for msg := range f.inChannel {
		f.forEachElement.Each(msg)
	}
}
