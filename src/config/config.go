package cfg

import (
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/labstack/gommon/log"
)

var (
	once sync.Once
	cfg config
)


type config struct {
	LogLevel string `env:"LOG_LEVEL" envDefaut:"info"`
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
		err := env.Parse(&cfg)
		if err !=nil{
			panic("config not set correctly")
		}

	})
	return cfg
}