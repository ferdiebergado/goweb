package handler

import (
	"database/sql"
	"net/http"

	"github.com/ferdiebergado/goexpress"
	"github.com/ferdiebergado/goweb/internal/config"
	"github.com/ferdiebergado/goweb/internal/pkg/security"
	"github.com/ferdiebergado/goweb/internal/repository"
	"github.com/ferdiebergado/goweb/internal/service"
	"github.com/go-playground/validator/v10"
)

type App struct {
	cfg       *config.Config
	db        *sql.DB
	router    *goexpress.Router
	validater *validator.Validate
	template  *Template
	hasher    security.Hasher
}

func NewApp(cfg *config.Config, db *sql.DB, r *goexpress.Router,
	v *validator.Validate, t *Template, h security.Hasher) *App {
	app := &App{
		cfg:       cfg,
		db:        db,
		router:    r,
		validater: v,
		template:  t,
		hasher:    h,
	}
	app.SetupMiddlewares()
	return app
}

func (a *App) Router() *goexpress.Router {
	return a.router
}

func (a *App) SetupMiddlewares() {
	a.router.Use(goexpress.RecoverFromPanic)
	a.router.Use(goexpress.LogRequest)
}

func (a *App) SetupRoutes() {
	if a.cfg.App.Env == "development" {
		const prefix = "/assets/"
		a.router.Handle("GET "+prefix, http.StripPrefix(prefix, http.FileServer(http.Dir("web/assets/"))))
	}

	repo := repository.NewRepository(a.db)
	svc := service.NewService(repo, a.hasher)

	htmlHandler := NewHandler(svc, a.template)
	apiHandler := NewAPIHandler(*svc)

	mountRoutes(a.router, htmlHandler)
	mountAPIRoutes(a.router, apiHandler, a.validater)
}
