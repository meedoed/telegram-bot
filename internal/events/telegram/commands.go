package telegram

import (
	"errors"
	"github.com/meedoed/telegram-bot/internal/lib/e"
	"github.com/meedoed/telegram-bot/internal/storage"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	// add page: http://...
	// rnd page: /rnd
	// help: /help
	// start: /start: hi + help

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sengHelp(chatID)
	case StartCmd:
		return p.sengHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) error {
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return e.Wrap("can't do command: save page", err)
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return e.Wrap("can't do command: save page", err)
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return e.Wrap("can't do command: save page", err)
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) error {
	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPage) {
		return e.Wrap("can't do command: can't send random", err)
	}
	if errors.Is(err, storage.ErrNoSavedPage) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return e.Wrap("can't do command: can't send random", err)
	}

	return p.storage.Remove(page)
}

func (p *Processor) sengHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sengHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
