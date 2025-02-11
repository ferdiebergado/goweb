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
}

func NewApp(cfg *config.Config, db *sql.DB, r *goexpress.Router, v *validator.Validate) *App {
	return &App{
		cfg:       cfg,
		db:        db,
		router:    r,
		validater: v,
	}
}

func (a *App) Router() *goexpress.Router {
	return a.router
}

func (a *App) SetupRoutes() {
	a.router.Use(goexpress.RecoverFromPanic)
	a.router.Use(goexpress.LogRequest)

	if a.cfg.App.Env == "development" {
		// TODO: extract prefix to const
		a.router.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets/"))))
	}

	repo := repository.NewRepository(a.db)
	baseService := service.NewService(repo)
	baseHandler := NewBaseHandler(baseService)
	mountRoutes(a.router, baseHandler)

	tmpl := NewTemplate(a.cfg.Template)
	baseHTMLHandler := NewBaseHTMLHandler(tmpl)
	mountBaseHTMLRoutes(a.router, baseHTMLHandler)

	userRepo := repository.NewUserRepository(a.db)
	hasher := &security.Argon2Hasher{}
	userService := service.NewUserService(userRepo, hasher)
	userHandler := NewUserHandler(userService, a.validater)
	mountUserRoutes(a.router, userHandler)
}
