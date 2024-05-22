package mapmark

import (
	"database/sql"
	"echoFramework/internal/config"
	"echoFramework/internal/handlers/mapmark"
	"echoFramework/internal/repository/mapmarkpg"
	mmservice "echoFramework/internal/service/mapmark"
	"errors"
	"fmt"
	trmsql "github.com/avito-tech/go-transaction-manager/drivers/sql/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"log/slog"
	"net/http"
	"os"
)

const DBDriverName = "sqlite3"

type App struct {
	conf   *config.Config
	logger *slog.Logger
	e      *echo.Echo
	db     *sql.DB
}

func New() *App {

	conf, err := config.New()
	if err != nil {
		log.Fatalf("can't init config: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := sql.Open(DBDriverName, conf.DBPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	mmRepo := mapmarkpg.New(db, trmsql.DefaultCtxGetter)
	mmService := mmservice.New(mmRepo)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	mapMarkGr := e.Group("/marks")
	mapmark.SetRouter(mmService, mapMarkGr)

	app := &App{
		conf:   conf,
		logger: logger,
		e:      e,
		db:     db,
	}

	return app
}

func (a *App) Run() {

	a.logger.Info("starting app...")

	go func() {
		if err := a.e.Start(fmt.Sprintf("%s:%s", a.conf.Host, a.conf.Port)); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				a.logger.Info("server stopped")
				return
			}
			a.logger.Error("can't start app", err)
		}
	}()

	a.logger.Info("app started")
}

func (a *App) Stop() {
	a.logger.Info("stopping app...")

	err := a.e.Shutdown(nil)
	if err != nil {
		a.logger.Error("can't shutdown app", err)
	}

	a.logger.Info("close db connection...")
	err = a.db.Close()
	if err != nil {
		a.logger.Error("can't close db connection", err)
	}

	a.logger.Info("app stopped")
}
