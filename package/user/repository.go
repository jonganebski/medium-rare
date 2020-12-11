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
	FindUserByEmail(user *model.User) (*model.User, error)
	FindUserByID(userID primitive.ObjectID) (*model.User, error)
	FindUsers(userIDs *[]primitive.ObjectID) (*[]model.User, error)
	InsertUser(user *model.User) (*primitive.ObjectID, error)
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
	filter := bson.D{{Key: "email", Value: user.Email}}
	singleResult := r.Collection.FindOne(context.Background(), filter)
	if err := singleResult.Err(); err != nil {
		return nil, err
	}
	singleResult.Decode(user)
	return user, nil
}

func (r *repository) FindUserByID(userID primitive.ObjectID) (*model.User, error) {
	user := new(model.User)
	filter := bson.D{{Key: "_id", Value: userID}}
	singleResult := r.Collection.FindOne(context.Background(), filter)
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
