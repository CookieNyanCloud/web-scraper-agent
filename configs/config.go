package configs

import (
	"flag"
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
			println(err.Error())
			return &Conf{}
		}
	}
	chat, _ := strconv.ParseInt(os.Getenv("CHAT"), 10, 64)
	return &Conf{
		os.Getenv("NOREGURL"),
		os.Getenv("MINNRURL"),
		os.Getenv("STARTMIN"),
		os.Getenv("TOKEN_A"),
		os.Getenv("URL"),
		chat,
	}
}
