package configs

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Conf struct {
	NoRegURL string
	MinNRURL string
	StartMin string
	Token    string
	URL      string
	Chat     int64
	NKOURL   string
	NKOBody  string
	FizURL   string
	Lasts *Lasts
}

type Lasts struct {
	LastNum int `json:"last_num"`
	LastNRNKO int `json:"last_nrnko"`
	NkoAll int `json:"nko_all"`
	Line int `json:"line"`
	LastFiz int `json:"last_fiz"`
	LastnoReg string `json:"lastno_reg"`
}

func InitConf() *Conf {
	var local bool
	flag.BoolVar(&local, "local", false, "хост")
	flag.Parse()
	return envVar(local)
}

func envVar(local bool) *Conf {
	if local {
		err := godotenv.Load(".env")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	chatID := os.Getenv("CHAT")
	chat, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	lasts,err:= parseLasts()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return &Conf{
		os.Getenv("NOREGURL"),
		os.Getenv("MINNRURL"),
		os.Getenv("STARTMIN"),
		os.Getenv("TOKEN_A"),
		os.Getenv("URL"),
		chat,
		os.Getenv("NKOURL"),
		os.Getenv("NKOBODY"),
		os.Getenv("FIZ_URL"),
		lasts,
	}
}

func parseLasts() (*Lasts,error){
	file, err := os.ReadFile("configs"+string(os.PathSeparator)+"lasts.json")
	if err != nil {
		return nil,err
	}
	out:= &Lasts{}
	err = json.Unmarshal(file, out)
	return out,err
}

func SaveLasts(cfg *Conf) error{
	err := os.Remove("configs"+string(os.PathSeparator)+"lasts.json")
	if err != nil {
		return err
	}
	file, err := os.Create("configs"+string(os.PathSeparator)+"lasts.json")
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(cfg.Lasts)
	if err != nil {
		return err
	}
	_, err = file.Write(bytes)
	if err != nil {
		return err
	}
	err = file.Close()
	return err
}