package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"

	"github.com/mwheatley3/ww/server/personal/api/api"
	"github.com/mwheatley3/ww/server/pg"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	// AuthTokenBytes is the length of the authentication token
	AuthTokenBytes = 32
)

func generateAuthToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), err
}

func hashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

// AuthenticateUser returns the *api.User corresponding to the email/password combination. If a user
// is not found, an error is returned.
func (db *Db) AuthenticateUser(email, password string) (*api.User, error) {
	user, err := db.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)) != nil {
		return nil, api.ErrInvalidPassword
	}

	return user, nil
}

// CreateUser creates a new User
func (db *Db) CreateUser(email, password string, sysAdmin bool) (*api.User, error) {
	var (
		err            error
		hashedPassword []byte
		v              dbUser
	)

	if hashedPassword, err = hashPassword([]byte(password)); err != nil {
		return nil, err
	}
	authToken, err := generateAuthToken(AuthTokenBytes)
	if err != nil {
		return nil, err
	}
	p := pg.NewParams(uuid.NewV4(), email, hashedPassword, sysAdmin, authToken)
	q := `INSERT INTO users (id, email, hashed_password, system_admin, auth_token) VALUES ($1, $2, $3, $4, $5) RETURNING ` + userCols
	err = db.db.Get(&v, q, p)
	if err != nil {
		return nil, db.unknownErr(err)
	}

	return v.User, nil
}

// UpdateUser updates a user. The password is only updated if the password field is not empty.
func (db *Db) UpdateUser(userID uuid.UUID, email, password string, sysAdmin bool) (*api.User, error) {
	var (
		v dbUser
		p = pg.NewParams(userID, email, sysAdmin)
		q = `UPDATE users SET email = $2, system_admin = $3, updated_at = now() WHERE id = $1 AND deleted_at IS NULL RETURNING ` + userCols
	)

	if password != "" {
		hashedPassword, err := hashPassword([]byte(password))
		if err != nil {
			return nil, err
		}
		authToken, err := generateAuthToken(AuthTokenBytes)
		if err != nil {
			return nil, err
		}
		p = pg.NewParams(userID, email, hashedPassword, sysAdmin, authToken)
		q = `UPDATE users SET email = $2, hashed_password = $3, system_admin = $4, auth_token = $5, updated_at = now() WHERE id = $1 AND deleted_at IS NULL RETURNING ` + userCols
	}

	err := db.db.Get(&v, q, p)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api.ErrInvalidUserID
		}
		return nil, db.unknownErr(err)
	}

	return v.User, nil
}

// DeleteUser deletes a user
func (db *Db) DeleteUser(userID uuid.UUID) (*api.User, error) {
	var (
		v dbUser
		p = pg.NewParams(userID)
		q = `UPDATE users SET updated_at = now(), deleted_at = now() WHERE id = $1 AND deleted_at IS NULL RETURNING ` + userCols
	)

	err := db.db.Get(&v, q, p)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api.ErrInvalidUserID
		}
		return nil, db.unknownErr(err)
	}

	return v.User, nil
}

// FindUserByID finds user by id
func (db *Db) FindUserByID(userID uuid.UUID) (*api.User, error) {
	var (
		v      dbUser
		params = pg.NewParams(userID)
		q      = `SELECT ` + userCols + ` FROM users WHERE id = $1 AND deleted_at IS NULL`
	)

	if err := db.db.Get(&v, q, params); err != nil {
		if err == sql.ErrNoRows {
			return nil, api.ErrInvalidUserID
		}

		return nil, db.unknownErr(err)
	}

	return v.User, nil
}

// FindUserByEmail finds user by email
func (db *Db) FindUserByEmail(email string) (*api.User, error) {
	var (
		v      dbUser
		params = pg.NewParams(email)
		q      = `SELECT ` + userCols + ` FROM users WHERE email = $1 AND deleted_at IS NULL`
	)

	if err := db.db.Get(&v, q, params); err != nil {
		if err == sql.ErrNoRows {
			return nil, api.ErrInvalidEmail
		}

		return nil, db.unknownErr(err)
	}

	return v.User, nil
}
