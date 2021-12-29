package webScraper

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/CookieNyanCloud/web-scraper-agent/configs"
	"github.com/gocolly/colly"
)

type Scraper struct {
	noRegURL string
	minNRURL string
	startMin string
	url      string
	lastNum  int
	dif      int
}

type IScraper interface {
	CheckNoReg() []string
	Check() bool
	Find() string
	GetLast() string
	GetLastNR() (string, error)
}

func NewScraper(conf *configs.Conf) IScraper {
	return &Scraper{
		noRegURL: conf.NoRegURL,
		minNRURL: conf.MinNRURL,
		startMin: conf.StartMin,
		url:      conf.URL,
		lastNum:  102,
		dif:      0,
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
			return ""
		}
	}
	fmt.Println(text)
	return text
}

func (s *Scraper) GetLast() string {
	var text, query string
	coll := colly.NewCollector()
	query = fmt.Sprint("tr:last-child td:nth-child(2)")
	coll.OnHTML("table tbody", func(e *colly.HTMLElement) {
		text = e.DOM.Find(query).Text() + "\n"
	})
	err := coll.Visit(s.url)
	if err != nil {
		log.Printf("err visiting %s: %v", s.url, err)
		return ""
	}
	fmt.Println(text)
	return text
}

func (s *Scraper) GetLastNR() (string, error) {
	var URL string
	coll := colly.NewCollector()
	coll.OnHTML("p a", func(e *colly.HTMLElement) {
		URL = e.Attr("href")
	})
	err := coll.Visit(s.minNRURL)
	if err != nil {
		log.Printf("err visiting %s: %v", s.url, err)
		return "", err
	}
	fmt.Println("URL ", s.startMin+URL)
	err = DownloadFile("noreg.xlsx", s.startMin+URL)
	return "", err
}

func (s *Scraper) CheckNoReg() []string {

	return nil
}

func DownloadFile(filepath, url string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:95.0) Gecko/20100101 Firefox/95.0")
	req.Header.Set("name", "value")
	req.Header.Set("name", "value")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}
