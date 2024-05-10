package telegram

import (
	"errors"
	"telebot/clients/telegram"
	"telebot/events"
	"telebot/lib/e"
	"telebot/storage/files"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage files.Storage
}

type Meta struct {
	ChatId   int
	Username string
}

var (
	ErrorUnknownEventType = errors.New("unknown event type")
	ErrorUnknownMeta      = errors.New("unknown meta")
)

func New(client *telegram.Client, storage files.Storage) *Processor {
	return &Processor{
		tg:      client,
		offset:  0,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.WrapError("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))
	for _, update := range updates {
		res = append(res, event(update))
	}

	p.offset = updates[len(updates)-1].Id + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	case events.Unknown:
		return e.WrapError("Unable to process, event is unknown", ErrorUnknownEventType)
	default:
		return e.WrapError("Unable to process, event type is not in list", ErrorUnknownEventType)
	}
}

func event(u telegram.Update) events.Event {
	res := events.Event{
		Type: fetchType(u),
		Text: fetchText(u),
	}

	if res.Type == events.Message {
		res.Meta = Meta{
			ChatId:   u.Message.Chat.Id,
			Username: u.Message.From.Username,
		}
	}

	return res
}

func fetchText(u telegram.Update) string {
	if u.Message == nil {
		return ""
	}

	return u.Message.Text
}

func fetchType(u telegram.Update) events.EventType {
	if u.Message == nil {
		return events.Unknown
	}

	return events.Message
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return err
	}

	return p.ProcessCommand(event.Text, meta.ChatId, meta.Username)
}

func meta(event events.Event) (Meta, error) {
	meta, ok := event.Meta.(Meta) //assertion that we have valid Meta
	if !ok {
		return Meta{}, e.WrapError("Unable to fetch meta for event", ErrorUnknownMeta)
	}

	return meta, nil
}
