package migrations

import (
	m "github.com/mwheatley3/ww/server/migrations"
)

func init() {
	Migrations.Add("admin.001_users",
		func(db m.DB) error {
			return db.ExecMany(
				m.Q(`CREATE DOMAIN password_type as TEXT CHECK (VALUE IN ('bcrypt'))`),
				m.Q(`
				CREATE TABLE users (
					id UUID NOT NULL,
					email TEXT NOT NULL,
					hashed_password BYTEA NOT NULL,
					password_type password_type NOT NULL DEFAULT 'bcrypt',
					auth_token TEXT NOT NULL,
					system_admin BOOL NOT NULL DEFAULT FALSE,
					created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
					updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
					deleted_at TIMESTAMPTZ DEFAULT NULL,

					PRIMARY KEY (id)
				)
			`),
				m.Q(`CREATE UNIQUE INDEX ON users (email) WHERE deleted_at IS NULL`),
				m.Q(`CREATE UNIQUE INDEX ON users (auth_token) WHERE deleted_at IS NULL`),
			)
		},

		func(db m.DB) error {
			return db.ExecMany(
				m.Q(`DROP TABLE users`),
				m.Q(`DROP DOMAIN password_type`),
			)
		},
	)
}
