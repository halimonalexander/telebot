package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"
	"telebot/clients/telegram"
	"telebot/storage"
)

const (
	COMMAND_ADD_PAGE = "/add"
	COMMAND_RND      = "/rnd"
	COMMAND_HELP     = "/help"
	COMMAND_START    = "/start"
)

var (
	ErrorPageAlreadyExists = errors.New("page already exists")
)

func (p *Processor) ProcessCommand(command string, chatId int, username string) error {
	command = strings.TrimSpace(command)

	log.Printf("got new commands '%' from '%'", command, username)

	var incomingUrl string
	if isAddCmd(command) {
		incomingUrl = command
		command = COMMAND_ADD_PAGE
	}

	sendMessage := newMessageSender(chatId, p.tg)

	switch command {
	case COMMAND_RND:
		page, err := p.getRandomPage(username)
		if err != nil {
			if errors.Is(err, storage.ErrorNoSavedPages) {
				return p.tg.SendMessage(chatId, MsgNxPage)
			} else {
				log.Println(err.Error())
				return p.tg.SendMessage(chatId, MsgFetchException)
			}
		}

		return p.tg.SendMessage(chatId, page.URL)
	case COMMAND_HELP:
		//return p.tg.SendMessage(chatId, MsgHelp)
		return sendMessage(MsgHelp)
	case COMMAND_START:
		return p.tg.SendMessage(chatId, MsgHello)
	case COMMAND_ADD_PAGE:
		err := p.savePage(incomingUrl, username)
		if err != nil {
			if errors.Is(err, ErrorPageAlreadyExists) {
				return p.tg.SendMessage(chatId, MsgPageAlreadyExists)
			} else {
				log.Println(err.Error())
				return p.tg.SendMessage(chatId, MsgStoreException)
			}
		}
		return p.tg.SendMessage(chatId, MsgSaved)
	default:
		return p.tg.SendMessage(chatId, MsgUnknownCommand)
	}
}

func (p *Processor) savePage(pageUrl string, username string) error {
	page := &storage.Page{
		URL:      pageUrl,
		UserName: username,
	}

	exists, err := p.storage.Exists(page)
	if err != nil {
		return err
	}

	if exists == true {
		return ErrorPageAlreadyExists
	}

	err = p.storage.Save(page)
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) getRandomPage(username string) (*storage.Page, error) {
	randomPage, err := p.storage.PickRandom(username)
	if err != nil {
		return nil, err
	}

	err = p.storage.Remove(randomPage)
	if err != nil {
		log.Println("Warning: unable to remove page")
	}

	return randomPage, nil
}

// Closure example
func newMessageSender(chatId int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatId, msg)
	}
}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
