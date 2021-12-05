package webScraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
)

type Scraper struct {
	url     string
	lastNum int
	dif     int
}

type IScraper interface {
	Check() bool
	Find() string
}

func NewScraper(url string) IScraper {
	return &Scraper{
		url:     url,
		lastNum: 102,
		dif:     0,
	}
}

func (s *Scraper) Check() bool {
	coll := colly.NewCollector()
	coll.AllowURLRevisit = true
	var numStr string
	coll.OnHTML("table tbody", func(e *colly.HTMLElement) {
		numStr = e.DOM.Find("tr:last-child td:nth-child(1)").Text()
	})
	err := coll.Visit(s.url)
	if err != nil {
		log.Printf("err visiting %s: %v", s.url, err)
	}
	if numStr == "" {
		return false
	}
	sf := strings.TrimSuffix(numStr, ".")
	numF, err := strconv.Atoi(sf)
	if err != nil {
		log.Printf("err getting last number: %v", err)
	}
	num := numF
	fmt.Println(num)
	if num > s.lastNum {
		fmt.Println(num, s.lastNum)
		s.dif = num - s.lastNum
		s.lastNum = num
		return true
	}
	return false
}

func (s *Scraper) Find() string {
	var text, query string
	for i := s.lastNum - s.dif + 2; i < s.lastNum+2; i++ {
		coll := colly.NewCollector()
		query = fmt.Sprintf("tr:nth-child(%d) td:nth-child(2)", i)
		coll.OnHTML("table tbody", func(e *colly.HTMLElement) {
			text += strconv.Itoa(i-1) + ")" + e.DOM.Find(query).Text() + "\n"
		})
		err := coll.Visit(s.url)
		if err != nil {
			log.Printf("err visiting %s: %v", s.url, err)
		}
	}
	fmt.Println(text)
	return text
}
