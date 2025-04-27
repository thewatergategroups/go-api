package entrypoints

import (
	cfg "go-api/src/config"
	"go-api/src/endpoints"
	"log"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

// API

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func server(cmd *cobra.Command, args []string){
	e := echo.New()
	e.Use(middleware.RequestID())  // 📌 Add unique ID to all logs/errors early
	e.Use(middleware.Logger())     // 📝 Log every request (with RequestID)
	e.Use(middleware.Recover())    // 🛑 Catch panics before they crash the server
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20))) // 🚦 Enforce before any work is done
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{Timeout: 10 * time.Second})) // ⏱ Enforce max duration
	e.Use(middleware.CSRF())       // 🛡 Security: CSRF protection
	e.Use(middleware.CORS())       // 🌍 Cross-origin access
	e.Use(middleware.Secure())     // 🔐 Security headers
	// Register validator
	e.Validator = &CustomValidator{validator: validator.New()}

	endpoints.RegisterGreetingsRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}

func GetServerCommand()*cobra.Command{
	return &cobra.Command{
		Use:   "server",
		Short: "Start the server",
		Run: server,
	}
}

// Migrations
func migrateUp(cmd *cobra.Command, args []string){
	if err := goose.Up(cfg.Db().DB, "src/migrations"); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
}
func migrateDown(cmd *cobra.Command, args []string){
	if err := goose.Down(cfg.Db().DB, "src/migrations"); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
}

func GetMigrateCmd()  *cobra.Command{
	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Manage database migrations",
	}
	var migrateUpCmd = &cobra.Command{
		Use:   "up",
		Short: "Apply all up migrations",
		Run: migrateUp,
	}
	
	var migrateDownCmd = &cobra.Command{
		Use:   "down",
		Short: "Rollback last migration",
		Run: migrateDown,
	}
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	return migrateCmd
}


