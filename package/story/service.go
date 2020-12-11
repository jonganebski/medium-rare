package story

import (
	"home/jonganebski/github/medium-rare/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Service is an interface from which our api module can access our repository of all our models
type Service interface {
	CreateStory(story *model.Story) (*primitive.ObjectID, error)
	FindStoryByID(storyID primitive.ObjectID) (*model.Story, error)
	FindStories(storyIDs *[]primitive.ObjectID) (*[]model.Story, error)
	FindRecentStories() (*[]model.Story, error)
	FindPickedStories() (*[]model.Story, error)
	FindPopularStories() (*[]model.Story, error)
	IncreaseViewCount(storyID primitive.ObjectID) (*model.Story, error)
	AddCommentID(storyID, commentID primitive.ObjectID) *fiber.Error
	UpdateLikedUserIDs(storyID, userID primitive.ObjectID, key string) *fiber.Error
	UpdateStoryBlock(storyID primitive.ObjectID, blocks *[]model.Block) *fiber.Error
	RemoveCommentID(storyID, commentID primitive.ObjectID) *fiber.Error
	RemoveStory(storyID primitive.ObjectID) *fiber.Error
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

func (s *service) RemoveStory(storyID primitive.ObjectID) *fiber.Error {
	return s.repository.DeleteStory(storyID)
}

func (s *service) UpdateStoryBlock(storyID primitive.ObjectID, blocks *[]model.Block) *fiber.Error {
	return s.repository.UpdateStoryBlock(storyID, blocks)
}

func (s *service) CreateStory(story *model.Story) (*primitive.ObjectID, error) {
	return s.repository.InsertStory(story)
}

func (s *service) UpdateLikedUserIDs(storyID, userID primitive.ObjectID, key string) *fiber.Error {
	return s.repository.UpdateLikedUserIDs(storyID, userID, key)
}

func (s *service) AddCommentID(storyID, commentID primitive.ObjectID) *fiber.Error {
	return s.repository.UpdateCommentID(storyID, commentID, "$push")
}

func (s *service) RemoveCommentID(storyID, commentID primitive.ObjectID) *fiber.Error {
	return s.repository.UpdateCommentID(storyID, commentID, "$pull")
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
	return s.repository.UpdateViewCount(storyID)
}
