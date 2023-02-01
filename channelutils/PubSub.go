package channelutils

import "log"

type AcceptAllFilter[T any] struct{}

func (a *AcceptAllFilter[T]) Filter(in T) bool {
	return true
}

type PublishSubsrcibe[T any] struct {
	inChan        chan T
	subscriptions []*Subscription[T]
}

type Subscription[T any] struct {
	OutOperations[T]
	p      *PublishSubsrcibe[T]
	filter FilterInterface[T]
}

// func (s *Subscription[T]) Out() <-chan T {
// return s.out
// }

func (s *Subscription[T]) CancelSubscription() {
	s.p.unsubscribe(s)
}

type FilterInterface[T any] interface {
	Filter(in T) bool
}

func (p *PublishSubsrcibe[T]) Publish(toPublish T) {
	p.inChan <- toPublish
}

func (p *PublishSubsrcibe[T]) Subscribe() *Subscription[T] {
	return p.SubscribeWithFilter(&AcceptAllFilter[T]{})
}

func (p *PublishSubsrcibe[T]) SubscribeWithFilter(filter FilterInterface[T]) *Subscription[T] {
	log.Println("subscribe")
	ret := &Subscription[T]{*NewOutOperations[T](), p, filter}
	p.subscriptions = append(p.subscriptions, ret)
	return ret
}

func (p *PublishSubsrcibe[T]) unsubscribe(s *Subscription[T]) {
	for i, activeSubscription := range p.subscriptions {
		if activeSubscription == s {
			close(s.outChannel)
			p.subscriptions = append(p.subscriptions[:i], p.subscriptions[i+1:]...)
			return
		}
	}
}

func (p *PublishSubsrcibe[T]) listen() {
	log.Println("listen")
	for inMsg := range p.inChan {
		for _, subscription := range p.subscriptions {
			if subscription.filter.Filter(inMsg) {
				subscription.outChannel <- inMsg
			}
		}
	}
}

func (p *PublishSubsrcibe[T]) Fill(in chan T) {
	p.inChan = in
	p.subscriptions = []*Subscription[T]{}
	go p.listen()
}

func NewPublishSubscribe[T any](in chan T) *PublishSubsrcibe[T] {
	ret := &PublishSubsrcibe[T]{in, make([]*Subscription[T], 0)}
	go ret.listen()
	return ret
}
