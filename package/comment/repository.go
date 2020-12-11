package comment

import (
	"context"
	"home/jonganebski/github/medium-rare/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//Repository interface allows us to access the CRUD Operations in mongo here.
type Repository interface {
	InsertComment(comment *model.Comment) (*model.Comment, error)
	FindComment(commentID primitive.ObjectID) (*model.Comment, error)
	FindComments(commnetIDs *[]primitive.ObjectID) (*[]model.Comment, error)
	DeleteComment(commentID primitive.ObjectID) *fiber.Error
	DeleteComments(commentIDs *[]primitive.ObjectID) *fiber.Error
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

func (r *repository) DeleteComments(commentIDs *[]primitive.ObjectID) *fiber.Error {
	f := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: commentIDs}}}}
	_, err := r.Collection.DeleteMany(context.Background(), f)
	if err != nil {
		return fiber.NewError(500, "Failed to delete")
	}
	return nil
}

func (r *repository) DeleteComment(commentID primitive.ObjectID) *fiber.Error {
	f := bson.D{{Key: "_id", Value: commentID}}
	deleteResult, err := r.Collection.DeleteOne(context.Background(), f)
	if err != nil {
		return fiber.NewError(500, "Failed to update")
	}
	if deleteResult.DeletedCount == 0 {
		return fiber.NewError(404, "Comment not found")
	}
	return nil
}

func (r *repository) InsertComment(comment *model.Comment) (*model.Comment, error) {
	insertOneResult, err := r.Collection.InsertOne(context.Background(), comment)
	if err != nil {
		return nil, err
	}
	f := bson.D{{Key: "_id", Value: insertOneResult.InsertedID}}
	singleResult := r.Collection.FindOne(context.Background(), f)
	if err := singleResult.Err(); err != nil {
		return nil, err
	}
	singleResult.Decode(comment)
	return comment, nil

}

func (r *repository) FindComment(commentID primitive.ObjectID) (*model.Comment, error) {
	comment := new(model.Comment)
	f := bson.D{{Key: "_id", Value: commentID}}
	singleResult := r.Collection.FindOne(context.Background(), f)
	if err := singleResult.Err(); err != nil {
		return nil, err
	}
	singleResult.Decode(comment)
	return comment, nil
}

func (r *repository) FindComments(commentIDs *[]primitive.ObjectID) (*[]model.Comment, error) {
	comments := make([]model.Comment, 0)
	f := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: commentIDs}}}}
	c, err := r.Collection.Find(context.Background(), f)
	if err != nil {
		return nil, err
	}
	if err = c.All(context.Background(), &comments); err != nil {
		return nil, err
	}
	return &comments, nil
}
