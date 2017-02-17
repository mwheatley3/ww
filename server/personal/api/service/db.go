package service

import (
	"github.com/Sirupsen/logrus"
	"github.com/mwheatley3/ww/server/personal/api/api"
	"github.com/mwheatley3/ww/server/personal/api/db"
	"github.com/satori/go.uuid"
)

// New returns a new database-backed service
func New(l *logrus.Logger, db *db.Db) api.Service {
	return &apiService{
		logger: l,
		db:     db,
	}
}

type apiService struct {
	logger *logrus.Logger
	db     *db.Db
}

func (s *apiService) Init() error {
	return s.db.Init()
}

func (s *apiService) AuthenticateUser(email, password string) (*api.User, error) {
	return s.db.AuthenticateUser(email, password)
}

func (s *apiService) CreateUser(email, password string) (*api.User, error) {
	return s.db.CreateUser(email, password, false)
}

func (s *apiService) GetUser(userID uuid.UUID) (*api.User, error) {
	return s.db.FindUserByID(userID)
}

func (s *apiService) GetUserByEmail(email string) (*api.User, error) {
	return s.db.FindUserByEmail(email)
}

func (s *apiService) UpdateUser(userID uuid.UUID, email, password string) (*api.User, error) {
	return s.db.UpdateUser(userID, email, password, false)
}

func (s *apiService) DeleteUser(userID uuid.UUID) (*api.User, error) {
	return s.db.DeleteUser(userID)
}
