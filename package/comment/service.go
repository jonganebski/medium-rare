package comment

import (
	"home/jonganebski/github/medium-rare/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Service is an interface from which our api module can access our repository of all our models
type Service interface {
	CreateComment(comment *model.Comment) (*model.Comment, error)
	FindComment(commentID primitive.ObjectID) (*model.Comment, error)
	FindComments(commentIDs *[]primitive.ObjectID) (*[]model.Comment, error)
	RemoveComment(commentID primitive.ObjectID) error
	RemoveComments(commentIDs *[]primitive.ObjectID) error
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

func (s *service) CreateComment(comment *model.Comment) (*model.Comment, error) {
	return s.repository.InsertComment(comment)
}

func (s *service) FindComment(commentID primitive.ObjectID) (*model.Comment, error) {
	return s.repository.FindComment(commentID)
}

func (s *service) FindComments(commentIDs *[]primitive.ObjectID) (*[]model.Comment, error) {
	return s.repository.FindComments(commentIDs)
}

func (s *service) RemoveComment(commentID primitive.ObjectID) error {
	return s.repository.DeleteComment(commentID)
}

func (s *service) RemoveComments(commentIDs *[]primitive.ObjectID) error {
	return s.repository.DeleteComments(commentIDs)
}
