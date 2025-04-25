package ui

import (
	"github.com/francescarpi/mytime/internal/config"
	"github.com/francescarpi/mytime/internal/repository"
	"github.com/francescarpi/mytime/internal/service"
	"github.com/francescarpi/mytime/internal/service/redmine"
)

type Dependencies struct {
	Service *service.Service
	Redmine *redmine.Redmine
}

func InitDeps() *Dependencies {
	cfg := config.Load()
	repo := repository.NewSqliteRepository(cfg.DBUrl)
	service := &service.Service{Repo: repo}
	redmine := redmine.NewRedmine(service)

	return &Dependencies{
		Service: service,
		Redmine: redmine,
	}
}
