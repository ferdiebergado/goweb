package main

import (
	"database/sql"
	"net/http"

	"github.com/ferdiebergado/goexpress"
	"github.com/ferdiebergado/goweb/internal/config"
	"github.com/ferdiebergado/goweb/internal/handler"
	"github.com/ferdiebergado/goweb/internal/repository"
	"github.com/ferdiebergado/goweb/internal/service"
)

type App struct {
	cfg    *config.Config
	db     *sql.DB
	router *goexpress.Router
}

func NewApp(cfg *config.Config, db *sql.DB, r *goexpress.Router) *App {
	return &App{
		cfg:    cfg,
		db:     db,
		router: r,
	}
}

func (a *App) Router() *goexpress.Router {
	return a.router
}

func (a *App) SetupRoutes() {
	a.router.Use(goexpress.RecoverFromPanic)
	a.router.Use(goexpress.LogRequest)

	if a.cfg.App.Env == "development" {
		a.router.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets/"))))
	}

	repo := repository.NewRepository(a.db)
	service := service.NewService(repo)
	baseHandler := handler.NewBaseHandler(service)
	mountRoutes(a.router, baseHandler)

	tmpl := handler.NewTemplate(a.cfg.Template)
	baseHTMLHandler := handler.NewBaseHTMLHandler(tmpl)
	mountBaseHTMLRoutes(a.router, baseHTMLHandler)
}
