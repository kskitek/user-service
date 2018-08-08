package user

import "gitlab.com/kskitek/arecar/user-service/http_boundary"

func NewHandler() http_boundary.Handler {
	return &handler{
		s: NewService(),
	}
}

func NewService() Service {
	return &crud{
		dao: NewDao(),
	}
}

func NewDao() Dao {
	return &InMemDao{}
}
