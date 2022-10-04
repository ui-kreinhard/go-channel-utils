package channelutils

type Map[T any, U any] struct {
	OutOperations[U]
	inChannel <-chan T
	mapper    Mapper[T, U]
}

type Mapper[T any, U any] interface {
	Map(T) U
}

func NewMap[T any, U any](inChannel <-chan T, mapper Mapper[T, U]) *Map[T, U] {
	ret := &Map[T, U]{*NewOutOperations[U](), inChannel, mapper}
	go ret.listen()
	return ret
}

func (m *Map[T, U]) listen() {
	for msg := range m.inChannel {
		u := m.mapper.Map(msg)
		m.outChannel <- u
	}
}
