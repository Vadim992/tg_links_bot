package telegram

import (
	"context"
	"errors"
	"links_tg-bot/lib/e"
	"links_tg-bot/storage"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatId int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command: %s from %s\n", text, username)
	// commands
	// add page: http://...
	// get random page: /rnd
	// help for people (how this bot works): /help
	// start work with bot: /start: hi + help

	if isAddCmd(text) {
		//TODO: addPage()
		return p.savePage(text, chatId, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatId, username)
	case HelpCmd:
		return p.sendHelp(chatId)
	case StartCmd:
		return p.sendHello(chatId)
	default:
		return p.tg.SendMessages(chatId, msgUnknownCommand)
	}

}

func (p *Processor) savePage(pageUrl string, chatId int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()
	page := &storage.Page{
		URL:      pageUrl,
		UserName: username,
	}

	//isExist, err := p.storage.IsExist(page)

	isExist, err := p.storage.IsExist(context.Background(), page)

	if err != nil {
		return err
	}

	if isExist {
		return p.tg.SendMessages(chatId, msgAlreadyExist)
	}

	//if err := p.storage.Save(page); err != nil {
	//	return err
	//}
	if err := p.storage.Save(context.Background(), page); err != nil {
		return err
	}

	if err := p.tg.SendMessages(chatId, msgSaved); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendRandom(chatId int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: send random page", err) }()

	//page, err := p.storage.PickRandom(username)

	page, err := p.storage.PickRandom(context.Background(), username)

	if err != nil && !errors.Is(err, storage.ErrorNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrorNoSavedPages) {
		return p.tg.SendMessages(chatId, msgNoSavedPages)
	}

	if err := p.tg.SendMessages(chatId, page.URL); err != nil {
		return err
	}

	//return p.storage.Remove(page)

	return p.storage.Remove(context.Background(), page)
}

func (p *Processor) sendHelp(chatId int) error {
	return p.tg.SendMessages(chatId, msgHelp)
}

func (p *Processor) sendHello(chatId int) error {
	return p.tg.SendMessages(chatId, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
