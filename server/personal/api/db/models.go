package db

import (
	"github.com/mwheatley3/ww/server/personal/api/api"
	"github.com/mwheatley3/ww/server/pg"
)

var (
	userCols = pg.Columns(&dbUser{}).String()
)

type dbUser struct {
	*api.User
}

func (v *dbUser) MapProxy() pg.MapProxy {
	if v.User == nil {
		v.User = &api.User{}
	}

	return pg.MapProxy{
		"users.id":              &v.ID,
		"users.email":           &v.Email,
		"users.hashed_password": &v.HashedPassword,
		"users.password_type":   &v.PasswordType,
		"users.auth_token":      &v.AuthToken,
		"users.system_admin":    &v.SystemAdmin,
		"users.created_at":      &v.CreatedAt,
		"users.updated_at":      &v.UpdatedAt,
		"users.deleted_at":      &v.DeletedAt,
	}
}
