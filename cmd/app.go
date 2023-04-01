package main

import (
	"context"
	"fmt"
	chat2 "github.com/Salavei/golang_websockets/internal/chat"
	"github.com/Salavei/golang_websockets/internal/chat/db"
	"github.com/Salavei/golang_websockets/internal/chat/handlers"
	"github.com/Salavei/golang_websockets/internal/config"
	"github.com/Salavei/golang_websockets/pkg/client/mongodb"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

func main() {
	cfg := config.GetConfig()
	ChatClient, err := mongodb.NewClient(context.Background(), cfg.Storage)
	if err != nil {
		log.Fatal(err)
	}

	storage := db.NewStorage(ChatClient, cfg.Storage.MongoDB.Collection)

	start(cfg, storage)
}

func start(cfg *config.Config, storage chat2.Storage) {

	webSocketConnect := chat.NewServer(storage)

	http.HandleFunc("/", config.ServeFiles)

	http.Handle("/chat", websocket.Handler(webSocketConnect.HandleWSChat))
	http.Handle("/online", websocket.Handler(webSocketConnect.HandleWSOnline))

	addr := fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	fmt.Printf("listening server on: %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
