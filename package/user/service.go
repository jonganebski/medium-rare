package user

import (
	"home/jonganebski/github/medium-rare/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Service is an interface from which our api module can access our repository of all our models
type Service interface {
	FindUserByEmail(user *model.User) (*model.User, error)
	FindUserByID(userID primitive.ObjectID) (*model.User, error)
	FindUsers(userIDs *[]primitive.ObjectID) (*[]model.User, error)
	CreateUser(user *model.User) (*primitive.ObjectID, error)
}

type service struct {
	repository Repository
}

//NewService is used to create a single instance of the service
func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) FindUsers(userIDs *[]primitive.ObjectID) (*[]model.User, error) {
	return s.repository.FindUsers(userIDs)
}

func (s *service) FindUserByEmail(user *model.User) (*model.User, error) {
	return s.repository.FindUserByEmail(user)
}

func (s *service) FindUserByID(userID primitive.ObjectID) (*model.User, error) {
	return s.repository.FindUserByID(userID)
}

func (s *service) CreateUser(user *model.User) (*primitive.ObjectID, error) {
	return s.repository.InsertUser(user)
}
