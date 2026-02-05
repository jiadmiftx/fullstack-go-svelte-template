package http_middleware

import (
	"net/http"

	"sekolah-madrasah/pkg/rbac"
)

var Router = []rbac.RouterRoles{
	{
		Endpoint: "/v1/organization/member/invite",
		Detail: []rbac.RouterDetail{
			{
				Roles:  []rbac.Role{rbac.Member, rbac.Owner, rbac.SuperAdmin},
				Method: http.MethodPost,
			},
			{
				Roles:  []rbac.Role{rbac.SuperAdmin},
				Method: http.MethodGet,
			},
		},
		SubRouter: []rbac.RouterRoles{
			{
				Endpoint: "/settings/:kegiatan_id",
				Detail: []rbac.RouterDetail{
					{Roles: []rbac.Role{rbac.SuperAdmin}, Method: http.MethodPut},
					{Roles: []rbac.Role{rbac.Member, rbac.Owner, rbac.SuperAdmin}, Method: http.MethodGet},
				},
			},
		},
	},
}
