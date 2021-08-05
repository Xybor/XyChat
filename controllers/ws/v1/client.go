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
func CreateWSClient(conn *websocket.Conn) wsClient {
	wsc := wsClient{
		conn:        conn,
		send:        make(chan interface{}),
		ReadHandler: func(message []byte) error { return nil },
	}

	go wsc.readFlow()
	go wsc.writeFlow()

	return wsc
}

func (wsc *wsClient) readFlow() {
	defer wsc.conn.Close()

	for {
		_, message, err := wsc.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
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
	wsc.CloseHandler()
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
			log.Println(err)
			break
		}
	}
}

// WriteJSON send a message to wsClient.writeFlow via wsClient.send, the
// message should be a struct or map[string]interface{}
func (wsc *wsClient) WriteJSON(data interface{}) {
	wsc.send <- data
}

// Close closes the connection and wsClient.send channel.
func (wsc *wsClient) Close() {
	close(wsc.send)
	wsc.conn.Close()
}
