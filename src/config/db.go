package cfg

import (
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var dbOnce sync.Once
var db *sqlx.DB

func Db() *sqlx.DB {
	dbOnce.Do(func(){
		url := fmt.Sprintf(
			"postgres://%s:%s@%s/%s?&search_path=%s",
			Cfg().Db.Username,
			Cfg().Secrets.PostgresPassword,
			Cfg().Db.Address,
			Cfg().Db.DbName,
			Cfg().Db.Schema,
		)
		var err error
		db, err = sqlx.Open(Cfg().Db.Driver, url)
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}
		db.SetMaxOpenConns(Cfg().Db.MaxOpenConns)                
		db.SetMaxIdleConns(Cfg().Db.MaxIdleConns)                 
		db.SetConnMaxLifetime(time.Duration(Cfg().Db.MaxConnLifetimeMins) * time.Minute) 
	
		if err := db.Ping(); err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}

		query:= fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`,Cfg().Db.Schema)
		_, err = db.Exec(query)
		if err != nil {
			log.Fatalf("failed to ensure schema exists: %v", err)
		}
	})
	return db
}
