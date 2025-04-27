package main

import (
	cfg "go-api/src/config"
	"go-api/src/entrypoints"
	_ "go-api/src/migrations"
	"log"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)


func init(){
	logger:= slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level:cfg.GetLogLevel(cfg.Cfg().LogLevel) },
		),
	)
	slog.SetDefault(logger)
	

}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "go-api",
		Short: "Go Api Template",
		Long:  "A CLI tool to run migrations or start the server",
	}
	rootCmd.AddCommand(entrypoints.GetMigrateCmd())
	rootCmd.AddCommand(entrypoints.GetServerCommand())
	
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}