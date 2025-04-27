package cfg

import (
	"encoding/json"
	"log/slog"
	"os"
	"sync"

	"github.com/caarlos0/env/v11"
)

var (
	once sync.Once
	cfg config
)

type config struct {
	LogLevel string `json:"log_level"`
	Cache string `json:"cache"`
	Redis redisConfig `json:"redis"`
	Db dbConfig `json:"db"`
	Secrets secrets
}

type redisConfig struct {
	Address string `json:"address"` // host and port ( localhost:6379 )
	Db int `json:"db"`
}

type dbConfig struct {
	Driver string `json:"driver"`
	Username string `json:"username"` 
	Address string `json:"address"`  // host and port ( localhost:6379 )
	DbName string `json:"db_name"`
	Schema string `json:"schema"`
	MaxOpenConns int `json:"max_open_conns"`
	MaxIdleConns int `json:"max_idle_conns"`
	MaxConnLifetimeMins int `json:"max_conn_lifetime_mins"`
}

type secrets struct {
	RedisPassword string `env:"REDIS_PASSWORD"`
	PostgresPassword string `env:"PG_PASSWORD"`
}

func GetLogLevel(logLevel string) slog.Level{
	level := slog.LevelInfo
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
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