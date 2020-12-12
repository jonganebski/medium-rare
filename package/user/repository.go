package user

import (
	"context"
	"fmt"
	"home/jonganebski/github/medium-rare/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository interface allows us to access the CRUD Operations in mongo here.
type Repository interface {
	InsertUser(user *model.User) (*primitive.ObjectID, error)
	FindUserByID(userID primitive.ObjectID) (*model.User, error)
	FindUserByEmail(user *model.User) (*model.User, error)
	FindUsers(userIDs *[]primitive.ObjectID) (*[]model.User, error)
	UpdateUserBookmark(userID, storyID primitive.ObjectID, key string) error
	UpdateCommentID(userID, commentID primitive.ObjectID, key string) error
	UpdateFollowingID(subjectUserID, targetUserID primitive.ObjectID, key string) error
	UpdateFollowerID(subjectUserID, targetUserID primitive.ObjectID, key string) error
	UpdateLikedStoryIDs(userID, storyID primitive.ObjectID, key string) error
	UpdateStoryIDs(userID, storyID primitive.ObjectID, key string) error
	UpdateUserDetails(userID primitive.ObjectID, field, value string) error
	RemoveManyLikedStoryIDs(storyID primitive.ObjectID) error
	RemoveManySavedStoryIDs(storyID primitive.ObjectID) error
	RemoveManyUsersCommentID(commentID primitive.ObjectID) error
	DeleteUser(userID primitive.ObjectID) error
}

type repository struct {
	Collection *mongo.Collection
}

//NewRepo is the single instance repo that is being created.
func NewRepo(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}

func (r *repository) DeleteUser(userID primitive.ObjectID) error {
	f := bson.D{{Key: "_id", Value: userID}}
	deleteResult, err := r.Collection.DeleteOne(context.Background(), f)
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount == 0 {
		return err
	}
	return nil
}

func (r *repository) RemoveManySavedStoryIDs(storyID primitive.ObjectID) error {
	f := bson.D{{Key: "savedStoryIds", Value: bson.D{{Key: "$elemMatch", Value: bson.D{{Key: "$eq", Value: storyID}}}}}}
	u := bson.D{{Key: "$pull", Value: bson.D{{Key: "savedStoryIds", Value: storyID}}}}
	_, err := r.Collection.UpdateMany(context.Background(), f, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) RemoveManyLikedStoryIDs(storyID primitive.ObjectID) error {
	f := bson.D{{Key: "likedStoryIds", Value: bson.D{{Key: "$elemMatch", Value: bson.D{{Key: "$eq", Value: storyID}}}}}}
	u := bson.D{{Key: "$pull", Value: bson.D{{Key: "likedStoryIds", Value: storyID}}}}
	_, err := r.Collection.UpdateMany(context.Background(), f, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateUserDetails(userID primitive.ObjectID, field, value string) error {
	f := bson.D{{Key: "_id", Value: userID}}
	u := bson.D{{Key: "$set", Value: bson.D{{Key: field, Value: value}}}}
	_, err := r.Collection.UpdateOne(context.Background(), f, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateStoryIDs(userID, storyID primitive.ObjectID, key string) error {
	f := bson.D{{Key: "_id", Value: userID}}
	u := bson.D{{Key: key, Value: bson.D{{Key: "storyIds", Value: storyID}}}}
	_, err := r.Collection.UpdateOne(context.Background(), f, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateLikedStoryIDs(userID, storyID primitive.ObjectID, key string) error {
	f := bson.D{{Key: "_id", Value: userID}}
	u := bson.D{{Key: key, Value: bson.D{{Key: "likedStoryIds", Value: storyID}}}}
	_, err := r.Collection.UpdateOne(context.Background(), f, u)
	if err != nil {
		return err
	}
	return nil
}

func (r repository) RemoveManyUsersCommentID(commentID primitive.ObjectID) error {
	f := bson.D{{Key: "commentIds", Value: bson.D{{Key: "$elemMatch", Value: bson.D{{Key: "$eq", Value: commentID}}}}}}
	u := bson.D{{Key: "$pull", Value: bson.D{{Key: "commentIds", Value: commentID}}}}
	_, err := r.Collection.UpdateOne(context.Background(), f, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateFollowingID(subjectUserID, targetUserID primitive.ObjectID, key string) error {
	f := bson.D{{Key: "_id", Value: subjectUserID}}
	u := bson.D{{Key: key, Value: bson.D{{Key: "followingIds", Value: targetUserID}}}}
	_, err := r.Collection.UpdateOne(context.Background(), f, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateFollowerID(subjectUserID, targetUserID primitive.ObjectID, key string) error {
	f := bson.D{{Key: "_id", Value: subjectUserID}}
	u := bson.D{{Key: key, Value: bson.D{{Key: "followerIds", Value: targetUserID}}}}
	_, err := r.Collection.UpdateOne(context.Background(), f, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateCommentID(userID, commentID primitive.ObjectID, key string) error {
	f := bson.D{{Key: "_id", Value: userID}}
	u := bson.D{{Key: key, Value: bson.D{{Key: "commentIds", Value: commentID}}}}
	_, err := r.Collection.UpdateOne(context.Background(), f, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateUserBookmark(userID, storyID primitive.ObjectID, key string) error {
	f := bson.D{{Key: "_id", Value: userID}}
	u := bson.D{{Key: key, Value: bson.D{{Key: "savedStoryIds", Value: storyID}}}}
	_, err := r.Collection.UpdateOne(context.Background(), f, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindUsers(userIDs *[]primitive.ObjectID) (*[]model.User, error) {
	users := make([]model.User, 0)
	f := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: userIDs}}}}
	c, err := r.Collection.Find(context.Background(), f)
	if err != nil {
		return nil, err
	}
	if err = c.All(context.Background(), &users); err != nil {
		return nil, err
	}
	return &users, nil
}

func (r *repository) FindUserByEmail(user *model.User) (*model.User, error) {
	f := bson.D{{Key: "email", Value: user.Email}}
	singleResult := r.Collection.FindOne(context.Background(), f)
	if err := singleResult.Err(); err != nil {
		return nil, err
	}
	singleResult.Decode(user)
	return user, nil
}

func (r *repository) FindUserByID(userID primitive.ObjectID) (*model.User, error) {
	user := new(model.User)
	f := bson.D{{Key: "_id", Value: userID}}
	singleResult := r.Collection.FindOne(context.Background(), f)
	if err := singleResult.Err(); err != nil {
		return nil, err
	}
	singleResult.Decode(user)
	return user, nil
}

func (r *repository) InsertUser(user *model.User) (*primitive.ObjectID, error) {
	insertOneResult, err := r.Collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}
	userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", insertOneResult.InsertedID))
	if err != nil {
		return nil, err
	}
	return &userOID, nil
}
