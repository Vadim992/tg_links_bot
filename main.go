package main

import (
	"context"
	"flag"
	tgClient "links_tg-bot/clients/telegram"
	event_consumer "links_tg-bot/consumer/event-consumer"
	"links_tg-bot/events/telegram"
	"links_tg-bot/storage/sqlite"
	"log"
	"strings"
)

const (
	tgHostBot = "api.telegram.org"
	//storagePath = "storage"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	// file storage
	//s := files.NewStorage(storagePath)

	// database storage

	s, err := sqlite.NewStorageDB(sqliteStoragePath)

	if err != nil {
		log.Fatalf("can't connect to storage: %e", err)
	}

	if err := s.InitDB(context.TODO()); err != nil {
		log.Fatalf("can't connect to storage: %e", err)
	}

	eventsProcessor := telegram.NewProcessor(
		tgClient.NewClient(tgHostBot, mustToken()),
		s,
	)
	log.Print("service started")

	consumer := event_consumer.NewConsumer(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}

}

func mustToken() string {
	var token string
	flag.StringVar(
		&token,
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	token = strings.TrimSpace(token)

	if token == "" {
		log.Fatal("token is not specified")
	}
	return token
}
