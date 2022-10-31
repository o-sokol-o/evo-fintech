package main

import (
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/o-sokol-o/evo-fintech/internal/config"
	"github.com/o-sokol-o/evo-fintech/internal/repositories"
	"github.com/o-sokol-o/evo-fintech/internal/services"
	"github.com/o-sokol-o/evo-fintech/internal/transport"
	"github.com/o-sokol-o/evo-fintech/pkg/database"
	"github.com/o-sokol-o/evo-fintech/pkg/server"
	"github.com/o-sokol-o/evo-fintech/pkg/signaler"
	"github.com/sirupsen/logrus"
)

// @title EVO Fintech API
// @version 1.0

func init() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
}

func main() {
	cfg, err := config.Init()
	if err != nil {
		logrus.Fatalln(err)
	}

	db, err := database.NewPostgresConnection(cfg.Postgres)
	if err != nil {
		logrus.Fatalln(err)
	}
	defer db.Close()

	repo := repositories.New(db)
	service := services.New(repo.RepoEVO, repo.RepoRemote)
	handler := transport.NewHandler(service.ServicesEVO, service.ServicesRemote)

	srv := server.New(cfg.HTTP.Port, handler.Init(cfg))

	go func() {
		if err = srv.Run(); err != nil {
			logrus.Errorf("error start server: %s", err.Error())
		}
	}()

	logrus.Println("rest-api started")

	signaler.Wait()

	if err = srv.Stop(); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
}
