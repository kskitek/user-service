package main

import (
	"net/http"
	"time"

	authsrv "github.com/kskitek/user-service/auth"
	"github.com/kskitek/user-service/event/redis"
	"github.com/kskitek/user-service/rest/auth"
	"github.com/kskitek/user-service/rest/user"
	"github.com/kskitek/user-service/server"
	ustore "github.com/kskitek/user-service/store/postgres/user"
	usersrv "github.com/kskitek/user-service/user"
	"github.com/sirupsen/logrus"
)

func main() {
	setupLogger()

	dao := ustore.NewPgDao()
	server.Server.AddRoutes(userHandler(dao).Routes())
	server.Server.AddRoutes(authHandler(dao).Routes())

	logrus.WithField("port", ":8080").Info("Starting")
	logrus.Fatal(http.ListenAndServe(":8080", server.Server.Handler()))
}

func setupLogger() {
	// TODO level from param
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: time.RFC3339})
}

func userHandler(dao usersrv.Dao) server.Handler {
	notifier, err := redis.NewNotifier()
	if err != nil {
		logrus.WithError(err).Fatal("unable to setup user handler")
	}
	service := usersrv.NewService(dao, notifier)
	return user.NewHandler(service)
}

func authHandler(dao usersrv.Dao) server.Handler {
	service := authsrv.NewService(dao)
	return auth.NewHandler(service)
}
