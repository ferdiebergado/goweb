package main

import (
	"database/sql"
	"net/http"

	"github.com/ferdiebergado/goexpress"
	"github.com/ferdiebergado/goweb/internal/config"
	"github.com/ferdiebergado/goweb/internal/handler"
	"github.com/ferdiebergado/goweb/internal/pkg/security"
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
	baseService := service.NewService(repo)
	baseHandler := handler.NewBaseHandler(baseService)
	mountRoutes(a.router, baseHandler)

	tmpl := handler.NewTemplate(a.cfg.Template)
	baseHTMLHandler := handler.NewBaseHTMLHandler(tmpl)
	mountBaseHTMLRoutes(a.router, baseHTMLHandler)

	userRepo := repository.NewUserRepository(a.db)
	hasher := &security.Argon2Hasher{}
	userService := service.NewUserService(userRepo, hasher)
	userHandler := handler.NewUserHandler(userService)
	mountUserRoutes(a.router, userHandler)
}
