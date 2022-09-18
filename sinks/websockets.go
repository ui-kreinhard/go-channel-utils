package sinks

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type ConnectionTabMap struct {
	connections map[string]*websocket.Conn
}

func NewConnectionsTabMap() *ConnectionTabMap {
	return &ConnectionTabMap{
		make(map[string]*websocket.Conn),
	}
}

func (c *ConnectionTabMap) Add(connectionId string, con *websocket.Conn) {
	c.connections[connectionId] = con
}

func (c *ConnectionTabMap) Remove(connectionId string) {
	delete(c.connections, connectionId)

	log.Println("removal", len(c.connections))
}

func (c *ConnectionTabMap) Length() int {
	return len(c.connections)
}

type WebsocketSink[T any] struct {
	inputChannel  <-chan T
	upgrader      websocket.Upgrader
	connections   ConnectionTabMap
	listenAddress string
	connCounter   int
}

func NewWebsocketSink[T any](inputChannel <-chan T, listenAddress string) *WebsocketSink[T] {
	ret := &WebsocketSink[T]{
		inputChannel,
		websocket.Upgrader{},
		*NewConnectionsTabMap(),
		listenAddress,
		0,
	}
	go ret.Listen()
	return ret
}

func (w *WebsocketSink[T]) shovelToWebsocket(conn *websocket.Conn, messageChan <-chan T, shutdown chan bool) {
	for {
		select {
		case frame := <-messageChan:
			err := conn.WriteJSON(frame)
			if err != nil {
				log.Println("err writing", err)
				shutdown <- true
				return
			}
		}
	}
}

func (ws *WebsocketSink[T]) tabWebsocket(w http.ResponseWriter, r *http.Request) {
	connectionId := ws.connCounter
	ws.connCounter++

	log.Println("new connection for", connectionId)

	if c, err := ws.upgrader.Upgrade(w, r, nil); err != nil {
		http.Error(w, fmt.Errorf("upgrade failed: %w", err).Error(), http.StatusInternalServerError)
		return
	} else {
		ws.connections.Add(strconv.Itoa(connectionId), c)

		shutdownChan := make(chan bool)

		go ws.shovelToWebsocket(c, ws.inputChannel, shutdownChan)

		<-shutdownChan
		ws.connections.Remove(strconv.Itoa(connectionId))
		defer c.Close()
		ws.idle()
	}
}

func (w *WebsocketSink[T]) idle() {
	log.Println("start idle")
	for {
		select {
		case <-w.inputChannel:
			if w.connections.Length() > 0 {
				log.Println("finish idle, connection occured", w.connections.Length())
				return
			}
		}
	}
}

func (w *WebsocketSink[T]) Listen() {
	log.Println("starting webserver on", w.listenAddress)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/connect", w.tabWebsocket)
	go w.idle()
	http.ListenAndServe(w.listenAddress, router)

}
