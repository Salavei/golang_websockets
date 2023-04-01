package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Salavei/golang_websockets/internal/chat"
	"golang.org/x/net/websocket"
	"io"
	"time"
)

var _ Server

type Server struct {
	conns   map[*websocket.Conn]bool
	conns2  map[*websocket.Conn]bool
	storage chat.Storage
}

func NewServer(storage chat.Storage) *Server {
	return &Server{
		conns:   make(map[*websocket.Conn]bool),
		conns2:  make(map[*websocket.Conn]bool),
		storage: storage,
	}
}

func (s *Server) HandleWSOnline(ws *websocket.Conn) {

	fmt.Println("new incoming connection online: ", ws.Request().RemoteAddr)
	s.conns2[ws] = true
	s.Online()
	s.onlineLoop(ws)

}

func (s *Server) onlineLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		_, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				delete(s.conns2, ws)
				s.Online()
				break
			}
			fmt.Println("read error: ", err)
			continue
		}
	}
}

func (s *Server) Online() {
	for ws := range s.conns2 {
		go func(ws *websocket.Conn) {
			sendMsg := fmt.Sprintf("%d", len(s.conns2))
			if _, err := ws.Write([]byte(sendMsg)); err != nil {
				fmt.Println("write error: ", err)
			}
		}(ws)
	}
}

func (s *Server) ShowAllMessages(msg []chat.Message, ws *websocket.Conn) {

	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Marshal error: ", err)
	}
	if _, err := ws.Write(data); err != nil {
		fmt.Println("write error: ", err)
	}
}

func (s *Server) HandleWSChat(ws *websocket.Conn) {

	fmt.Println("new incoming connection chat: ", ws.Request().RemoteAddr)
	s.conns[ws] = true
	data, err := s.storage.ShowMessage(context.Background())
	if err != nil {
		fmt.Errorf("messages dont show : %v ", err.Error())
	}
	s.ShowAllMessages(data, ws)

	s.readLoop(ws)

}

type Message struct {
	Type string        `json:"type"`
	Text string        `json:"message"`
	ID   string        `json:"user_id"`
	Date time.Duration `json:"date"`
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error: ", err)
			continue
		}
		msg := buf[:n]
		var msgStruct Message
		json.Unmarshal(msg, &msgStruct)

		ct := chat.Message{
			Text:   msgStruct.Text,
			UserID: msgStruct.ID,
			Date:   msgStruct.Date.Abs(),
		}
		s.storage.SendMessage(context.Background(), ct.UserID, ct)
		s.broadcast(msgStruct)
	}
}

func (s *Server) broadcast(msg Message) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			data, err := json.Marshal(msg)
			if err != nil {
				fmt.Println("Marshal error: ", err)
			}
			if _, err := ws.Write(data); err != nil {
				fmt.Println("write error: ", err)
			}
		}(ws)
	}
}
