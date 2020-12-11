package story

import (
	"home/jonganebski/github/medium-rare/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Service is an interface from which our api module can access our repository of all our models
type Service interface {
	FindRecentStories() (*[]model.Story, error)
	FindPickedStories() (*[]model.Story, error)
	FindPopularStories() (*[]model.Story, error)
	IncreaseViewCount(storyID primitive.ObjectID) (*model.Story, error)
	FindStoryByID(storyID primitive.ObjectID) (*model.Story, error)
	FindStories(storyIDs *[]primitive.ObjectID) (*[]model.Story, error)
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

func (s *service) FindStories(storyIDs *[]primitive.ObjectID) (*[]model.Story, error) {
	return s.repository.FindStories(storyIDs)
}

func (s *service) FindStoryByID(storyID primitive.ObjectID) (*model.Story, error) {
	return s.repository.FindStoryByID(storyID)
}

func (s *service) FindRecentStories() (*[]model.Story, error) {
	return s.repository.FindRecentStories()
}

func (s *service) FindPickedStories() (*[]model.Story, error) {
	return s.repository.FindPickedStories()
}

func (s *service) FindPopularStories() (*[]model.Story, error) {
	return s.repository.FindPopularStories()
}

func (s *service) IncreaseViewCount(storyID primitive.ObjectID) (*model.Story, error) {
	return s.repository.IncreaseViewCount(storyID)
}
