package main

import (
	"fmt"
	"strconv"
	"strings"
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
				if (time.Now().Hour()+3)%24 >= 10 || (time.Now().Hour()+3)%24 <= 1 {
					if scraper.Check() {
						s, _ := scraper.Find()
						for k := range users {
							msg := tgbotapi.NewMessage(k, "объявлены иноагентами:\n"+s)
							_, _ = bot.Send(msg)
							msgURL := tgbotapi.NewMessage(k, conf.URL)
							_, _ = bot.Send(msgURL)
						}
						msg := tgbotapi.NewMessage(conf.Chat, "объявлены иноагентами:\n"+s)
						_, _ = bot.Send(msg)
						msgURL := tgbotapi.NewMessage(conf.Chat, conf.URL)
						_, _ = bot.Send(msgURL)
						//toChannel := tgbotapi.NewMessageToChannel(conf.Chan, "объявлены иноагентами:\n"+s)
						//_, _ = bot.Send(toChannel)
					}
					noRegnko, err := scraper.CheckNoReg()
					if err != nil {
						fmt.Printf("err check no reg nko:%v", err)
					}
					if noRegnko {
						nko, err := scraper.FindNoRegNKO()
						if err != nil {
							fmt.Printf("err check no reg nko:%v", err)
						}
						for k := range users {
							msg := tgbotapi.NewMessage(k, "объявлены иноагентами:\n"+nko)
							_, _ = bot.Send(msg)
							msgURL := tgbotapi.NewMessage(k, conf.NoRegURL)
							_, _ = bot.Send(msgURL)
						}
						msg := tgbotapi.NewMessage(conf.Chat, "объявлены иноагентами:\n"+nko)
						_, _ = bot.Send(msg)
						msgURL := tgbotapi.NewMessage(conf.Chat, conf.NoRegURL)
						_, _ = bot.Send(msgURL)
					}
					lastNKO, all, err := scraper.GetLastNKO()
					if err != nil {
						fmt.Printf("err last nko:%v", err)
					}
					if lastNKO {
						for k := range users {
							msg := tgbotapi.NewMessage(k, "объявлены новые НКО:\n"+strconv.Itoa(all))
							_, _ = bot.Send(msg)
							msgURL := tgbotapi.NewMessage(k, conf.NKOURL)
							_, _ = bot.Send(msgURL)

						}
						msg := tgbotapi.NewMessage(conf.Chat, "объявлены новые НКО:\n"+strconv.Itoa(all))
						_, _ = bot.Send(msg)
						msgURL := tgbotapi.NewMessage(conf.Chat, conf.NKOURL)
						_, _ = bot.Send(msgURL)
					}

					if scraper.CheckZapr() {
						s := strings.Join(scraper.FindZapr(), "\n")
						for k := range users {
							msg := tgbotapi.NewMessage(k, "запрещены сайты:\n"+s)
							_, _ = bot.Send(msg)
							msgURL := tgbotapi.NewMessage(k, conf.URL)
							_, _ = bot.Send(msgURL)
						}
						msg := tgbotapi.NewMessage(conf.Chat, "запрещены сайты:\n"+s)
						_, _ = bot.Send(msg)
						msgURL := tgbotapi.NewMessage(conf.Chat, conf.URL)
						_, _ = bot.Send(msgURL)

					}
				}
			}
		}

	}()

	for update := range updates {

		if update.Message == nil {
			continue
		}

		if update.Message.Command() == "check" {

			lastNoReg, err := scraper.GetLastNR()
			if err != nil {
				fmt.Printf("no reg: %v", err)
			}
			lastSMI := scraper.GetLast()

			lastNKO, all, err := scraper.GetLastNKO()
			if err != nil {
				fmt.Printf("nko: %v", err)
			}
			lastZapr := scraper.GetLastZapr()
			if err != nil {
				fmt.Printf("сайт: %v", err)
			}
			textSMI := fmt.Sprintf("последний в списке иноагентов:%s\n\n", lastSMI)
			textNRNKO := fmt.Sprintf("последний в списке незарегестрированных НКО:%s\n", lastNoReg)
			textNKO := fmt.Sprintf("новые в  списке НКО:%t,%d\n", lastNKO, all)
			textZapr := fmt.Sprintf("новые запрещенные сайты:%s\n", lastZapr)
			allURL := fmt.Sprintf("%s\n%s\n%s\n%s\n", conf.URL, conf.NKOURL, conf.NoRegURL, conf.ZaprURL)
			fmt.Println(conf)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, textSMI+textNRNKO+textNKO+textZapr+allURL)
			_, _ = bot.Send(msg)
			continue
		}

		if update.Message.Command() == "time" {
			t1 := (time.Now().Hour()+3)%24 >= 16
			t2 := (time.Now().Hour()+3)%24 <= 1
			t3 := time.Now()
			text := fmt.Sprintf("t1>=16 %t\nt2<=1 %t\ntime %v\nall %v\n", t1, t2, t3, len(users))
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			_, _ = bot.Send(msg)
			continue
		}
		_, ok := users[update.Message.Chat.ID]
		if !ok {
			users[update.Message.Chat.ID] = struct{}{}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "слежу")
			_, _ = bot.Send(msg)
			continue
		}

	}

}
