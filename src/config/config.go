package cfg

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/labstack/gommon/log"
)

var (
	once sync.Once
	cfg config
)

type config struct {
	LogLevel string `json:"log_level"`
	Redis redis `json:"redis"`
	Secrets secrets
}

type redis struct {
	Enabled bool `json:"enabled"`
	Address string `json:"address"` // host and port
	Db int `json:"db"`
}

type secrets struct {
	RedisPassword string `env:"REDIS_PASSWORD"`
}

func GetLogLevel(logLevel string) log.Lvl{
	level := log.INFO
	switch logLevel {
	case "debug":
		level = log.DEBUG
	case "warn":
		level = log.WARN
	case "error":
		level = log.ERROR
	}
	return level
}

func Cfg() config{
	once.Do(func(){
		
		configFile, err := os.Open("config.json")
		if err != nil {
			panic(err.Error())
		}
		defer configFile.Close()
		if err:= json.NewDecoder(configFile).Decode(&cfg);err != nil {
			panic(err.Error())
		}
		
		if err := env.Parse(&cfg.Secrets); err !=nil{
			panic(err.Error())
		}

	})
	return cfg
}