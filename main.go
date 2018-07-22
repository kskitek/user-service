package main

import (
	"net/http"
	"gitlab.com/kskitek/arecar/user-service/http_boundary"
)

func main() {
	s := http_boundary.NewServer(getRoutes())
	http.ListenAndServe(":8080", s.Router)
}
