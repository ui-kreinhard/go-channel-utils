package channelutils

type Map[T any, U any] struct {
	inChannel  <-chan T
	outChannel chan U
	mapper     Mapper[T, U]
}

type Mapper[T any, U any] interface {
	Map(T) U
}

func NewMap[T any, U any](inChannel <-chan T, mapper Mapper[T, U]) *Map[T, U] {
	ret := &Map[T, U]{inChannel, make(chan U), mapper}
	go ret.listen()
	return ret
}

func (m *Map[T, U]) listen() {
	for msg := range m.inChannel {
		u := m.mapper.Map(msg)
		m.outChannel <- u
	}
}

func (m *Map[T, U]) Out() <-chan U {
	return m.outChannel
}

func (m *Map[T, U]) Filter(filerElement FilterElement[U]) *Filter[U] {
	return NewFilter(m.Out(), filerElement)
}

func (m *Map[T, U]) Mux() *ChannelMultiplexer[U] {
	return NewChannelMultiplexer(m.Out())
}

func (m *Map[T, U]) Each(eachElement ForEachElement[U]) *ForEach[U] {
	return NewForEach(m.Out(), eachElement)
}
