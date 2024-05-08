package telegram

import (
	"telebot/clients/telegram"
	"telebot/storage/files"
)

type Processor struct {
	tg     *telegram.Client
	offset int
	s      files.Storage
}

func New(client *telegram.Client, s files.Storage) Processor {
	return Processor{
		tg:     client,
		offset: 0,
		s:      s,
	}
}
