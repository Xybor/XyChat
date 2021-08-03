package v1

import (
	"log"

	"github.com/gorilla/websocket"
)

type wsClient struct {
	conn         *websocket.Conn
	send         chan interface{}
	ReadHandler  func(message []byte) error
	CloseHandler func()
}

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

func (wsc wsClient) readFlow() {
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

func (wsc wsClient) writeFlow() {
	defer wsc.conn.Close()

	for {
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

func (wsc wsClient) WriteJSON(data interface{}) {
	wsc.send <- data
}

func (wsc wsClient) Close() {
	close(wsc.send)
	wsc.conn.Close()
}
