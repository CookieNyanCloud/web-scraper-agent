package main

import (
	"github.com/CookieNyanCloud/web-scraper-agent/configs"
	"github.com/CookieNyanCloud/web-scraper-agent/sotatgbot"
	"github.com/CookieNyanCloud/web-scraper-agent/webScraper"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func main() {

	conf := configs.InitConf()
	users := make(map[int64]struct{}, 1)

	scraper := webScraper.NewScraper(conf.URL)
	bot, updates := sotatgbot.StartSotaBot(conf.Token)

	go func() {
		ticker := time.NewTicker(time.Minute)
		for {
			select {
			case <-ticker.C:
				if scraper.Check() {
					s := scraper.Find()
					for k, _ := range users {
						msg := tgbotapi.NewMessage(k, "объявлены иноагентами:\n"+s)
						_, _ = bot.Send(msg)
					}
				}
			}
		}

	}()

	for update := range updates {

		if update.Message == nil {
			continue
		} else {
			users[update.Message.Chat.ID] = struct{}{}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "слежу")
			_, _ = bot.Send(msg)
		}

	}

}
