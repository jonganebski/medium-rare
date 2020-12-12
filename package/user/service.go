package user

import (
	"home/jonganebski/github/medium-rare/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Service is an interface from which our api module can access our repository of all our models
type Service interface {
	CreateUser(user *model.User) (*primitive.ObjectID, error)
	FindUserByEmail(user *model.User) (*model.User, error)
	FindUserByID(userID primitive.ObjectID) (*model.User, error)
	FindUsers(userIDs *[]primitive.ObjectID) (*[]model.User, error)
	BookmarkStory(userID, storyID primitive.ObjectID) error
	DisbookmarkStory(userID, storyID primitive.ObjectID) error
	AddStoryID(userID, storyID primitive.ObjectID) error
	AddCommentID(userID, commentID primitive.ObjectID) error
	AddFollowerID(subjectUserID, targetUserID primitive.ObjectID) error
	AddFollowingID(subjectUserID, targetUserID primitive.ObjectID) error
	UpdateLikedStoryIDs(userID, storyID primitive.ObjectID, key string) error
	RemoveFollowerID(subjectUserID, targetUserID primitive.ObjectID) error
	RemoveFollowingID(subjectUserID, targetUserID primitive.ObjectID) error
	RemoveCommentID(userID, commentID primitive.ObjectID) error
	RemoveStoryID(userID, storyID primitive.ObjectID) error
	EditUsername(userID primitive.ObjectID, value string) error
	EditBio(userID primitive.ObjectID, value string) error
	EditAvatar(userID primitive.ObjectID, value string) error
	EditPassword(userID primitive.ObjectID, value string) error
	RemoveManyLikedStoryIDs(storyID primitive.ObjectID) error
	RemoveManySavedStoryIDs(storyID primitive.ObjectID) error
	RemoveManyUsersCommentID(commentID primitive.ObjectID) error
	RemoveAccount(userID primitive.ObjectID) error
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

func (s *service) RemoveManyUsersCommentID(commentID primitive.ObjectID) error {
	return s.repository.RemoveManyUsersCommentID(commentID)
}

func (s *service) RemoveAccount(userID primitive.ObjectID) error {
	return s.repository.DeleteUser(userID)
}

func (s *service) RemoveManyLikedStoryIDs(storyID primitive.ObjectID) error {
	return s.repository.RemoveManyLikedStoryIDs(storyID)
}

func (s *service) RemoveManySavedStoryIDs(storyID primitive.ObjectID) error {
	return s.repository.RemoveManySavedStoryIDs(storyID)
}

func (s *service) EditPassword(userID primitive.ObjectID, value string) error {
	return s.repository.UpdateUserDetails(userID, "password", value)
}

func (s *service) EditAvatar(userID primitive.ObjectID, value string) error {
	return s.repository.UpdateUserDetails(userID, "avatarUrl", value)
}

func (s *service) EditBio(userID primitive.ObjectID, value string) error {
	return s.repository.UpdateUserDetails(userID, "bio", value)
}

func (s *service) EditUsername(userID primitive.ObjectID, value string) error {
	return s.repository.UpdateUserDetails(userID, "username", value)
}

func (s *service) AddStoryID(userID, storyID primitive.ObjectID) error {
	return s.repository.UpdateStoryIDs(userID, storyID, "$push")
}

func (s *service) RemoveStoryID(userID, storyID primitive.ObjectID) error {
	return s.repository.UpdateStoryIDs(userID, storyID, "$pull")
}

func (s *service) UpdateLikedStoryIDs(userID, storyID primitive.ObjectID, key string) error {
	return s.repository.UpdateLikedStoryIDs(userID, storyID, key)
}

func (s *service) AddFollowingID(subjectUserID, targetUserID primitive.ObjectID) error {
	return s.repository.UpdateFollowingID(subjectUserID, targetUserID, "$push")
}

func (s *service) RemoveFollowingID(subjectUserID, targetUserID primitive.ObjectID) error {
	return s.repository.UpdateFollowingID(subjectUserID, targetUserID, "$pull")
}

func (s *service) AddFollowerID(subjectUserID, targetUserID primitive.ObjectID) error {
	return s.repository.UpdateFollowerID(subjectUserID, targetUserID, "$push")
}

func (s *service) RemoveFollowerID(subjectUserID, targetUserID primitive.ObjectID) error {
	return s.repository.UpdateFollowerID(subjectUserID, targetUserID, "$pull")
}

func (s *service) AddCommentID(userID, commentID primitive.ObjectID) error {
	return s.repository.UpdateCommentID(userID, commentID, "$push")
}

func (s *service) RemoveCommentID(userID, commentID primitive.ObjectID) error {
	return s.repository.UpdateCommentID(userID, commentID, "$pull")
}

func (s *service) DisbookmarkStory(userID, storyID primitive.ObjectID) error {
	return s.repository.UpdateUserBookmark(userID, storyID, "$pull")
}

func (s *service) BookmarkStory(userID, storyID primitive.ObjectID) error {
	return s.repository.UpdateUserBookmark(userID, storyID, "$push")
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
