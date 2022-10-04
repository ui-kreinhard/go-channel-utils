package channelutils

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Command string

const (
	Skip            Command = "s"
	Continue        Command = "c"
	SetStopeMode    Command = "sm"
	SetContinueMode Command = "cm"
)

type Mode string

const (
	StopOnMessage     Mode = "Stop"
	ContinueOnMessage Mode = "ContinueOnMessage"
)

type DebugWS[T any] struct {
	OutOperations[T]
	listenAddress           string
	inChannel               <-chan T
	websocketControlChannel chan Command
	websocketOutChannel     chan string
	upgrader                websocket.Upgrader
	mode                    Mode
	c                       *websocket.Conn
}

func NewDebugWS[T any](inChannel <-chan T, debugStartMode Mode, listenAddress string) *DebugWS[T] {
	ret := &DebugWS[T]{
		*NewOutOperations[T](),
		listenAddress,
		inChannel,
		make(chan Command),
		make(chan string),
		websocket.Upgrader{},
		debugStartMode,
		nil,
	}
	go ret.shovel()

	go ret.startWebserver()
	return ret
}

func (d *DebugWS[T]) shovel() {
	for msg := range d.inChannel {

		if d.mode == StopOnMessage {
			d.websocketOutChannel <- fmt.Sprint(msg)
			cmd := <-d.websocketControlChannel
			switch cmd {
			case Skip:
				log.Println("skipping output")
			case Continue:
				d.outChannel <- msg
			}
		} else if d.mode == ContinueOnMessage {
			log.Println("continue mode")
			d.websocketOutChannel <- fmt.Sprint(msg)
			d.outChannel <- msg
		}
	}
}

func (d *DebugWS[T]) Out() <-chan T {
	return d.outChannel
}

func (d *DebugWS[T]) shovelFromWebsocket(c *websocket.Conn) {
	for {
		var cmdMsg Command
		log.Println("waiting for cmd")
		err := c.ReadJSON(&cmdMsg)
		log.Println("got cmd", cmdMsg, err)
		switch cmdMsg {
		case SetStopeMode:
			d.mode = StopOnMessage
			c.WriteJSON("Currently in " + d.mode)
		case SetContinueMode:
			d.mode = ContinueOnMessage
			c.WriteJSON("Currently in " + d.mode)
			d.websocketControlChannel <- Continue
		case Continue:
			c.WriteJSON("Sending out")
			d.websocketControlChannel <- cmdMsg
		case Skip:
			c.WriteJSON("Skipping")
			d.websocketControlChannel <- cmdMsg
		default:
			d.printUsage(c)
		}

	}
}

func (d *DebugWS[T]) printUsage(c *websocket.Conn) {
	usage := "" +
		"c = continue to output\n" +
		"s = skip from output\n" +
		"sm = Stop on every msg(stopMode)\n" +
		"cm = Don't stop on every msg(continue mode)\n" +
		"enclose cmd in \""
	c.WriteJSON(usage)
}

func (d *DebugWS[T]) shovelToWebsocket(c *websocket.Conn) {
	for msg := range d.websocketOutChannel {
		c.WriteJSON(msg)
	}
}

func (d *DebugWS[T]) connect(w http.ResponseWriter, r *http.Request) {
	log.Println("ws connected")
	if c, err := d.upgrader.Upgrade(w, r, nil); err != nil {
		http.Error(w, fmt.Errorf("upgrade failed: %w", err).Error(), http.StatusInternalServerError)
		return
	} else {
		c.WriteJSON("Currently in " + d.mode)
		go d.shovelToWebsocket(c)
		d.shovelFromWebsocket(c)
		defer c.Close()
	}
}

func (d *DebugWS[T]) startWebserver() {
	log.Println("starting debug T listener on", d.listenAddress)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/connect", d.connect)
	http.ListenAndServe(d.listenAddress, router)
}
