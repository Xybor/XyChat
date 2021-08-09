package v1

import (
	"log"

	"github.com/gorilla/websocket"
)

// The wsClient is responsible for managing a websocket connection.
type wsClient struct {
	conn         *websocket.Conn
	send         chan interface{}
	ReadHandler  func(message []byte) error
	CloseHandler func()

	// Connection will be close if and only writeFlow is finished
	endOfWriteFlow chan bool
}

// Create a wsClient struct with the provided websocket connection.  It runs
// two goroutines for reading and writing flows.
//
// The wsClient.readFlow reads all messages from client and passes them into
// wsClient.ReadHandler function.
//
// The wsClient.writeFlow receives message from wsClient.send channel and sends
// it to the client.  The message passed to wsClient.send channel should be a
// struct or map[string]interface{}.
//
// Before wsClient.readFlow finishes, it call wsClient.CloseHandler function.
func CreateWSClient(conn *websocket.Conn) *wsClient {
	wsc := wsClient{
		conn:           conn,
		send:           make(chan interface{}),
		ReadHandler:    func(message []byte) error { return nil },
		CloseHandler:   func() {},
		endOfWriteFlow: make(chan bool),
	}

	go wsc.readFlow()
	go wsc.writeFlow()

	return &wsc
}

func (wsc *wsClient) readFlow() {
	defer wsc.conn.Close()

	for {
		_, message, err := wsc.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseNormalClosure,
			) {
				log.Println(err)
			}
			break
		}

		err = wsc.ReadHandler(message)
		if err != nil {
			log.Println(err)
			break
		}
	}
<<<<<<< HEAD
	//wsc.CloseHandler()
=======

	wsc.CloseHandler()
>>>>>>> 946e20c13db55d5c113aac5b51fa3a0ed0f8f59f
}

func (wsc *wsClient) writeFlow() {
	defer wsc.conn.Close()

	for {
		// The received message should be a struct or map.
		msg, ok := <-wsc.send
		if !ok {
			break
		}

		err := wsc.conn.WriteJSON(msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseNormalClosure,
			) {
				log.Println(err)
			}
			break
		}
	}

	wsc.endOfWriteFlow <- true
}

// WriteJSON send a message to wsClient.writeFlow via wsClient.send, the
// message should be a struct or map[string]interface{}
func (wsc *wsClient) WriteJSON(data interface{}) {
	wsc.send <- data
}

// Close closes the connection and wsClient.send channel.
func (wsc *wsClient) Close() {
	close(wsc.send)
	<-wsc.endOfWriteFlow
	wsc.conn.Close()
}
