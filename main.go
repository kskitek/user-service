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
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go/config"
)

func main() {
	setupLogger()
	setupTracer()

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

func setupTracer() {
	tracer, _, err := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: false,
		},
		ServiceName: "user-service",
	}.NewTracer()
	if err != nil {
		logrus.WithError(err).Error()
	}
	//tracer, _ := jaeger.NewTracer("user-service", jaeger.NewConstSampler(true), jaeger.NewInMemoryReporter())

	opentracing.InitGlobalTracer(tracer)
}
