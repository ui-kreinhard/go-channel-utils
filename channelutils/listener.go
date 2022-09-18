package channelutils

type Listener interface {
	Listen()
}

func Listen(listeners ...Listener) {
	for _, listener := range listeners {
		go listener.Listen()
	}
}
