package main

import (
	"fmt"
	"time"

	"github.com/CookieNyanCloud/web-scraper-agent/configs"
	"github.com/CookieNyanCloud/web-scraper-agent/sotatgbot"
	"github.com/CookieNyanCloud/web-scraper-agent/webScraper"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	conf := configs.InitConf()
	users := make(map[int64]struct{}, 1)

	scraper := webScraper.NewScraper(conf)
	bot, updates := sotatgbot.StartSotaBot(conf.Token)

	go func() {
		ticker := time.NewTicker(time.Minute)
		for {
			select {
			case <-ticker.C:
				if (time.Now().Hour()+3)%24 >= 15 || (time.Now().Hour()+3)%24 <= 1 {
					if scraper.Check() {
						s := scraper.Find()
						for k, _ := range users {
							msg := tgbotapi.NewMessage(k, "объявлены иноагентами:\n"+s)
							_, _ = bot.Send(msg)
						}
					}
				}
			}
		}

	}()

	for update := range updates {

		if update.Message == nil {
			continue
		} else if update.Message.Command() == "check" {
			noReg, err := scraper.GetLastNR()
			if err != nil {
				fmt.Printf("no reg: %v", err)
			}
			fmt.Println(noReg)
			last := scraper.GetLast()
			t1 := (time.Now().Hour()+3)%24 >= 16
			t2 := (time.Now().Hour()+3)%24 <= 1
			t3 := time.Now()
			text := fmt.Sprintf("последний в списке иноагент:%s\n%v>=16:%t\n%v<=1:%t\n(%v)", last, (time.Now().Hour()+3)%24, t1, (time.Now().Hour()+3)%24, t2, t3)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			_, _ = bot.Send(msg)
		} else {
			users[update.Message.Chat.ID] = struct{}{}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "слежу")
			_, _ = bot.Send(msg)
		}

	}

}
