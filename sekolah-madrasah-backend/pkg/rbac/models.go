package rbac

import "github.com/gin-gonic/gin"

type Role int64

const (
	Unauthorized Role = iota
	Member
	SuperAdmin
	Owner
	ExceptUnauthorized
)

type RouterDetail struct {
	Roles               []Role
	Tags                []string
	Method              string
	CustomParamsMessage string
	CustomMiddleware    []gin.HandlerFunc
	MustUrlParams       []string
	MustJsonBody        []string
	MustMultiPartBody   []string
}

type RouterRoles struct {
	Endpoint  string
	Detail    []RouterDetail
	SubRouter []RouterRoles
}

type Filter struct {
	Method string
	Role   *Role
}

type FindDataType interface {
	RouterRoles | RouterDetail
}
