package story

import (
	"context"
	"fmt"
	"home/jonganebski/github/medium-rare/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Repository interface allows us to access the CRUD Operations in mongo here.
type Repository interface {
	FindRecentStories() (*[]model.Story, error)
	FindPickedStories() (*[]model.Story, error)
	FindPopularStories() (*[]model.Story, error)
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

func (r *repository) FindRecentStories() (*[]model.Story, error) {
	stories := make([]model.Story, 0)
	o := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(30)
	f := bson.D{{Key: "isPublished", Value: true}, {Key: "editorsPick", Value: false}}
	c, err := r.Collection.Find(context.Background(), f, o)
	if err != nil {
		fmt.Println("Error at FindRecentStories")
		return nil, err
	}
	if err := c.All(context.Background(), &stories); err != nil {
		fmt.Println("Error at FindRecentStories")
		return nil, err
	}
	return &stories, nil
}

func (r *repository) FindPickedStories() (*[]model.Story, error) {
	stories := make([]model.Story, 0)
	o := options.Find().SetLimit(5)
	f := bson.D{{Key: "isPublished", Value: true}, {Key: "editorsPick", Value: true}}
	c, err := r.Collection.Find(context.Background(), f, o)
	if err != nil {
		fmt.Println("Error FindPickedStories")
		return nil, err
	}

	if err := c.All(context.Background(), &stories); err != nil {
		fmt.Println("error at FindPickedStories")
		return nil, err
	}
	return &stories, nil
}

func (r *repository) FindPopularStories() (*[]model.Story, error) {
	stories := make([]model.Story, 0)
	o := options.Find().SetSort(bson.D{{Key: "viewCount", Value: -1}}).SetLimit(5)
	f := bson.D{{Key: "isPublished", Value: true}}
	c, err := r.Collection.Find(context.Background(), f, o)
	if err != nil {
		fmt.Println("Error at FindPopularStories")
		return nil, err
	}
	if err := c.All(context.Background(), &stories); err != nil {
		fmt.Println("error at FindPopularStories")
		return nil, err
	}
	return &stories, nil
}
