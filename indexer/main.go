package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"w3todo-indexer/internal/api"
	"w3todo-indexer/internal/config"
	"w3todo-indexer/internal/db"
	"w3todo-indexer/internal/indexer"
)

func main() {
	cfg := config.FromEnv()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	database, err := db.New(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("БД: %v", err)
	}
	defer database.Close()

	idx, err := indexer.New(ctx, cfg, database)
	if err != nil {
		log.Fatalf("Индексатор: %v", err)
	}
	go idx.Start(ctx)

	server := api.New(database)
	log.Printf("API на :%s", cfg.Port)
	go http.ListenAndServe(":"+cfg.Port, server.Router())

	<-ctx.Done()
	log.Println("Завершение...")
}
