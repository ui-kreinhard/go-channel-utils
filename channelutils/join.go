package channelutils

import "log"

type Join[T any, U any, V any] struct {
	fast <-chan T
	slow <-chan U
	OutOperations[V]
	join JoinElement[T, U, V]

	inCurrent1 *T
	inCurrent2 *U
	onlyOnSlow bool
}

type JoinElement[T any, U any, V any] interface {
	Join(in1 T, in2 U) V
}

func NewJoin[T any, U any, V any](fast <-chan T, slow <-chan U, joinElement JoinElement[T, U, V]) *Join[T, U, V] {
	ret := &Join[T, U, V]{
		fast,
		slow,
		*NewOutOperations[V](),
		joinElement,
		nil,
		nil,
		false,
	}
	go ret.listen()
	return ret
}

func (j *Join[T, U, V]) SlowIsMaster() *Join[T, U, V] {
	j.onlyOnSlow = true
	return j
}

func (j *Join[T, U, V]) NoMaster() *Join[T, U, V] {
	j.onlyOnSlow = false
	return j
}

func (j *Join[T, U, V]) Out() <-chan V {
	return j.outChannel
}

func (j *Join[T, U, V]) listen() {
	go j.listenOnFast()
	go j.listenOnSlow()
}

func (j *Join[T, U, V]) listenOnFast() {
	for element1 := range j.fast {
		j.inCurrent1 = &element1
		if !j.onlyOnSlow {
			j.commit()
		}
	}
}

func (j *Join[T, U, V]) listenOnSlow() {
	for element2 := range j.slow {
		log.Println(j.inCurrent2)
		j.inCurrent2 = &element2
		j.commit()
	}
}

func (j *Join[T, U, V]) commit() {
	if j.inCurrent1 != nil && j.inCurrent2 != nil {
		in1 := *j.inCurrent1
		in2 := *j.inCurrent2

		j.inCurrent1 = nil
		j.inCurrent2 = nil

		j.outChannel <- j.join.Join(in1, in2)
	}
}
