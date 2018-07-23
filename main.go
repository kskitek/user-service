package main

import (
	"net/http"
	"gitlab.com/kskitek/arecar/user-service/http_boundary"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	s := http_boundary.NewServer(getRoutes())
	logrus.Fatal(http.ListenAndServe(":8080", s.Router))
}
