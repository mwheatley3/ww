package api

import "github.com/satori/go.uuid"

// Service is the main functional definition of the portal service
type Service interface {
	Init() error
	AuthenticateUser(email, password string) (*User, error)
	CreateUser(email, password string) (*User, error)
	GetUser(userID uuid.UUID) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(userID uuid.UUID, email, encryptedPassword string) (*User, error)
	DeleteUser(userID uuid.UUID) (*User, error)
}
