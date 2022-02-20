package main

import (
	"flag"
	tgClient "github.com/meedoed/telegram-bot/internal/clients/telegram"
	event_consumer "github.com/meedoed/telegram-bot/internal/consumer/event-consumer"
	"github.com/meedoed/telegram-bot/internal/events/telegram"
	"github.com/meedoed/telegram-bot/internal/storage/files"
	"log"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"token-bot-token",
		"",
		"token for acess to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
