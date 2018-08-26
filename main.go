package main

import (
	"net/http"
	"time"

	"github.com/kskitek/user-service/http_boundary"
	"github.com/sirupsen/logrus"
)

func main() {
	setupLogger()

	s := http_boundary.NewServer(getRoutes())
	logrus.WithField("port", ":8080").Info("Starting")
	logrus.Fatal(http.ListenAndServe(":8080", s.Router))
}

func setupLogger() {
	// TODO level from param
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: time.RFC3339})
}
