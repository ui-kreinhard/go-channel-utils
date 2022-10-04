package channelutils

type Filter[T any] struct {
	OutOperations[T]
	inChannel     <-chan T
	filterElement FilterElement[T]
}

func NewFilter[T any](inChannel <-chan T, filterElement FilterElement[T]) *Filter[T] {
	ret := &Filter[T]{
		*NewOutOperations[T](),
		inChannel,
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
