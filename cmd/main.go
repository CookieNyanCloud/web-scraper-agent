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
	scraper := webScraper.NewScraper(conf)
	bot, updates := sotatgbot.StartSotaBot(conf.Token)

	go func(conf *configs.Conf) {
		ticker := time.NewTicker(time.Minute)
		for {
			select {
			case <-ticker.C:
				if (time.Now().Hour()+3)%24 >= 14 || (time.Now().Hour()+3)%24 <= 1 {
					// smi
					if scraper.Check() {
						s, _ := scraper.Find()
						msg := tgbotapi.NewMessage(conf.Chat, "объявлены иноагентами:\n"+s+"\n"+conf.URL)
						_, _ = bot.Send(msg)
						err := configs.SaveLasts(conf)
						if err != nil {
							fmt.Printf("err SaveLasts:%v", err)
						}
					}

					// no reg
					noRegnko, err := scraper.CheckNoReg()
					if err != nil {
						fmt.Printf("err check no reg nko:%v", err)
					}
					if noRegnko {
						nko, err := scraper.FindNoRegNKO()
						if err != nil {
							fmt.Printf("err check no reg nko:%v", err)
						}
						msg := tgbotapi.NewMessage(conf.Chat, "объявлены иноагентами:\n"+nko+"\n"+conf.NoRegURL)
						_, _ = bot.Send(msg)
						err = configs.SaveLasts(conf)
						if err != nil {
							fmt.Printf("err SaveLasts:%v", err)
						}
					}

					// nko
					lastNKO, all, err := scraper.GetLastNKO()
					if err != nil {
						fmt.Printf("err last nko:%v", err)
					}
					if lastNKO {
						msg := tgbotapi.NewMessage(conf.Chat, "объявлены новые НКО:\n"+strconv.Itoa(all)+"\n"+conf.NKOURL)
						_, _ = bot.Send(msg)
						err := configs.SaveLasts(conf)
						if err != nil {
							fmt.Printf("err SaveLasts:%v", err)
						}
					}

					// fiz
					if scraper.CheckFiz() {
						s, _ := scraper.FindFiz()
						msg := tgbotapi.NewMessage(conf.Chat, "объявлены иноагентами физлицами:\n"+s+"\n"+conf.FizURL)
						_, _ = bot.Send(msg)
						err := configs.SaveLasts(conf)
						if err != nil {
							fmt.Printf("err SaveLasts:%v", err)
						}
					}
				}
			}
		}

	}(conf)

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
			lastFiz := scraper.GetLastFiz()
			if err != nil {
				fmt.Printf("физ: %v", err)
			}
			textSMI := fmt.Sprintf("последний в списке иноагентов:%s\n\n", lastSMI)
			textNRNKO := fmt.Sprintf("последний в списке незарегестрированных НКО:%s\n\n", lastNoReg)
			textNKO := fmt.Sprintf("новые в  списке НКО:%t,%d\n\n", lastNKO, all)
			textZapr := fmt.Sprintf("новые запрещенные физлица:%s\n\n", lastFiz)
			allURL := fmt.Sprintf("%s\n%s\n%s\n%s\n\n", conf.URL, conf.NKOURL, conf.NoRegURL, conf.FizURL)
			fmt.Println(conf)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, textSMI+textNRNKO+textNKO+textZapr+allURL)
			_, _ = bot.Send(msg)
			continue
		}

		if update.Message.Command() == "time" {
			t1 := (time.Now().Hour()+3)%24 >= 13
			t2 := (time.Now().Hour()+3)%24 <= 1
			t3 := time.Now()
			state := (time.Now().Hour()+3)%24 >= 14 || (time.Now().Hour()+3)%24 <= 1
			text := fmt.Sprintf("t1>=13 %t\nt2<=1 %t\ntime %v\n state %v", t1, t2, t3, state)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			_, _ = bot.Send(msg)
			continue
		}

		if update.Message.Command() == "setall" {
			i, err := strconv.Atoi(update.Message.CommandArguments())
			if err != nil {
				continue
			}
			scraper.SetAll(i)
			continue
		}
		if update.Message.Command() == "setline" {
			i, err := strconv.Atoi(update.Message.CommandArguments())
			if err != nil {
				continue
			}
			scraper.SetLine(i)
			continue
		}
		if update.Message.Command() == "getdata" {

			out := strings.Join(scraper.GetData(), ",")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, out)
			_, _ = bot.Send(msg)
			continue
		}

		if update.Message.Command() == "getLast" {

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%v",conf.Lasts))
			_, _ = bot.Send(msg)
			continue
		}

	}

}
