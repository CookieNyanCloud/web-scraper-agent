package webScraper

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/CookieNyanCloud/web-scraper-agent/configs"
	"github.com/gocolly/colly"
	"github.com/xuri/excelize/v2"
)

type Scraper struct {
	noRegURL  string
	minNRURL  string
	startMin  string
	url       string
	lastNum   int
	dif       int
	lastNRNKO int
	nkoURL    string
	nkoAll    int
	line      int
	nkoBody   string
	zaprURL   string
	lastZapr  string
	FizURL    string
	lastFiz   int
	difFiz    int
	lastnoReg string
}

type IScraper interface {
	SetLine(line int)
	SetAll(all int)
	getData() []string
	//media
	Check() bool
	Find() (string, int)
	GetLast() string
	//no reg nko
	CheckNoReg() (bool, error)
	FindNoRegNKO() (string, error)
	GetLastNR() (string, error)
	//nko
	GetLastNKO() (bool, int, error)
	//fiz
	CheckFiz() bool
	FindFiz() (string, int)
	GetLastFiz() string
}

func NewScraper(conf *configs.Conf) IScraper {
	return &Scraper{
		noRegURL:  conf.NoRegURL,
		minNRURL:  conf.MinNRURL,
		startMin:  conf.StartMin,
		url:       conf.URL,
		lastNum:   141,
		dif:       0,
		lastNRNKO: 9,
		nkoAll:    75,
		line:      405,
		nkoURL:    conf.NKOURL,
		nkoBody:   conf.NKOBody,
		lastFiz:   1,
		FizURL:    conf.FizURL,
		difFiz:    0,
		lastnoReg: "Инициативная группа ЛГБТ+ «Реверс»",
	}
}
func (s *Scraper) SetLine(line int) {
	s.line = line
}

func (s *Scraper) SetAll(all int) {
	s.nkoAll = all
}

func (s *Scraper) getData() []string {
	out := make([]string, 2)
	out = append(out, strconv.Itoa(s.line))
	out = append(out, strconv.Itoa(s.nkoAll))
	return out
}

//smi
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
	if num > s.lastNum {
		fmt.Println(num, s.lastNum)
		s.dif = num - s.lastNum
		s.lastNum = num
		return true
	}
	return false
}

func (s *Scraper) Find() (string, int) {
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
			return "", 0
		}
	}
	fmt.Println(text)
	return text, s.dif
}

func (s *Scraper) GetLast() string {
	var text, query string
	coll := colly.NewCollector()
	query = "tr:last-child td:nth-child(2)"
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

//no reg
func (s *Scraper) CheckNoReg() (bool, error) {
	last, err := s.GetLastNR()
	if err != nil {
		return false, err
	}
	if last == s.lastnoReg {
		return false, nil
	}
	s.lastnoReg = last
	return true, nil
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
	if err != nil {
		return "", err
	}
	f, err := excelize.OpenFile("noreg.xlsx")
	if err != nil {
		return "", err
	}
	rows, err := f.GetRows("Лист1")
	if err != nil {
		return "", err
	}
	err = os.Remove("noreg.xlsx")
	if err != nil {
		return "", err
	}
	return rows[len(rows)-1][1], nil
}

func (s *Scraper) FindNoRegNKO() (string, error) {

	err := DownloadFile("noreg.xlsx", s.noRegURL)
	if err != nil {
		return "", err
	}
	f, err := excelize.OpenFile("noreg.xlsx")
	if err != nil {
		return "", err
	}
	rows, err := f.GetRows("Лист1")
	if err != nil {
		return "", err
	}
	err = os.Remove("noreg.xlsx")
	if err != nil {
		return "", err
	}
	var out string

	for i := s.lastNRNKO; i < len(rows); i++ {
		out += fmt.Sprintf("%v%v\n", rows[i][0], rows[i][1])
	}
	s.lastNRNKO = len(rows)
	return out, nil
}

//nko
func (s *Scraper) GetLastNKO() (bool, int, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", s.nkoURL, strings.NewReader(s.nkoBody))
	if err != nil {
		return false, 0, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:96.0) Gecko/20100101 Firefox/96.0")
	resp, err := client.Do(req)
	if err != nil {
		return false, 0, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	if err != nil {
		return false, 0, err
	}

	line := 1
	for scanner.Scan() {
		if line != s.line {
			line++
			continue
		}
		check := fmt.Sprintf("[1&nbsp;-&nbsp;%d]", s.nkoAll)
		if !strings.Contains(scanner.Text(), check) {
			for i := s.nkoAll; i < 300; i++ {
				if strings.Contains(scanner.Text(), "[1&nbsp;-&nbsp;"+strconv.Itoa(i)+"]") {
					s.nkoAll = i
					break
				}
				if strings.Contains(scanner.Text(), "[1&nbsp;-&nbsp;"+strconv.Itoa(300-i)+"]") {
					s.nkoAll = 300 - i
					break
				}
			}
			return true, s.nkoAll, nil
		}
		break
	}
	if err := scanner.Err(); err != nil {
		return false, 0, err
	}
	return false, s.nkoAll, err

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

// fiz
func (s *Scraper) CheckFiz() bool {
	coll := colly.NewCollector()
	coll.AllowURLRevisit = true
	var numStr string
	coll.OnHTML("table tbody", func(e *colly.HTMLElement) {
		numStr = e.DOM.Find("tr:last-child td:nth-child(1)").Text()
	})
	err := coll.Visit(s.FizURL)
	if err != nil {
		log.Printf("err visiting %s: %v", s.FizURL, err)
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
	if num > s.lastFiz {
		fmt.Println(num, s.lastFiz)
		s.difFiz = num - s.lastFiz
		s.lastFiz = num
		return true
	}
	return false
}

func (s *Scraper) FindFiz() (string, int) {
	var text, query string
	for i := s.lastFiz - s.difFiz + 2; i < s.lastFiz+2; i++ {
		coll := colly.NewCollector()
		query = fmt.Sprintf("tr:nth-child(%d) td:nth-child(2)", i)
		coll.OnHTML("table tbody", func(e *colly.HTMLElement) {
			text += strconv.Itoa(i-1) + ")" + e.DOM.Find(query).Text() + "\n"
		})
		err := coll.Visit(s.FizURL)
		if err != nil {
			log.Printf("err visiting %s: %v", s.FizURL, err)
			return "", 0
		}
	}
	fmt.Println(text)
	return text, s.difFiz
}

func (s *Scraper) GetLastFiz() string {
	var text, query string
	coll := colly.NewCollector()
	query = "tr:last-child td:nth-child(2)"
	coll.OnHTML("table tbody", func(e *colly.HTMLElement) {
		text = e.DOM.Find(query).Text() + "\n"
	})
	err := coll.Visit(s.FizURL)
	if err != nil {
		log.Printf("err visiting %s: %v", s.FizURL, err)
		return ""
	}
	return text
}
