package ui

import (
	"github.com/francescarpi/mytime/internal/config"
	"github.com/francescarpi/mytime/internal/repository"
	"github.com/francescarpi/mytime/internal/service"
)

type Dependencies struct {
	Service *service.Service
}

func InitDeps() *Dependencies {
	cfg := config.Load()
	repo := repository.NewSqliteRepository(cfg.DBUrl)
	svc := &service.Service{Repo: repo}

	return &Dependencies{
		Service: svc,
	}
}
