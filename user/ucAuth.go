package user

import "gitlab.com/kskitek/arecar/user-service/http_boundary"

type Auth interface {
	Login(int, string) (bool, *http_boundary.ApiError)
}

type auth struct {

}