package main

import (
	"net/http"
	"gitlab.com/kskitek/arecar/user-service/http_boundary"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	setupLogger()

	s := http_boundary.NewServer(getRoutes())
	logrus.Fatal(http.ListenAndServe(":8080", s.Router))
}

func setupLogger() {
	// TODO level from param
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: time.RFC3339})
}
