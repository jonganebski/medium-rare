package story

import (
	"context"
	"fmt"
	"home/jonganebski/github/medium-rare/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Repository interface allows us to access the CRUD Operations in mongo here.
type Repository interface {
	InsertStory(story *model.Story) (*primitive.ObjectID, error)
	FindStoryByID(storyID primitive.ObjectID) (*model.Story, error)
	FindStories(storyIDs *[]primitive.ObjectID) (*[]model.Story, error)
	FindRecentStories(timestamp int64) (*[]model.Story, error)
	FindPickedStories() (*[]model.Story, error)
	FindPopularStories() (*[]model.Story, error)
	UpdateViewCount(storyID primitive.ObjectID) (*model.Story, error)
	UpdateCommentID(storyID, commentID primitive.ObjectID, key string) error
	UpdateLikedUserIDs(storyID, userID primitive.ObjectID, key string) error
	UpdatePickUnpick(storyID primitive.ObjectID, isPicked bool) error
	UpdatePublishStatus(storyID primitive.ObjectID, isPublished bool) error
	UpdateStoryBlock(storyID primitive.ObjectID, blocks *[]model.Block) error
	DeleteStory(storyID primitive.ObjectID) error
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

func (r *repository) UpdatePickUnpick(storyID primitive.ObjectID, isPicked bool) error {
	f := bson.D{{Key: "_id", Value: storyID}}
	u := bson.D{{Key: "$set", Value: bson.D{{Key: "editorsPick", Value: isPicked}}}}
	_, err := r.Collection.UpdateOne(context.Background(), f, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdatePublishStatus(storyID primitive.ObjectID, isPublished bool) error {
	f := bson.D{{Key: "_id", Value: storyID}}
	u := bson.D{{Key: "$set", Value: bson.D{{Key: "isPublished", Value: isPublished}}}}
	_, err := r.Collection.UpdateOne(context.Background(), f, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) DeleteStory(storyID primitive.ObjectID) error {
	f := bson.D{{Key: "_id", Value: storyID}}
	deleteResult, err := r.Collection.DeleteOne(context.Background(), f)
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount == 0 {
		return err
	}
	return nil
}

func (r *repository) UpdateStoryBlock(storyID primitive.ObjectID, blocks *[]model.Block) error {
	f := bson.D{{Key: "_id", Value: storyID}}
	u := bson.D{{Key: "$set", Value: bson.D{{Key: "blocks", Value: blocks}, {Key: "updatedAt", Value: time.Now().Unix()}}}}
	_, err := r.Collection.UpdateOne(context.Background(), f, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) InsertStory(story *model.Story) (*primitive.ObjectID, error) {
	insertOneResult, err := r.Collection.InsertOne(context.Background(), story)
	if err != nil {
		return nil, err
	}
	storyOID := insertOneResult.InsertedID.(primitive.ObjectID)
	return &storyOID, nil
}

func (r *repository) UpdateLikedUserIDs(storyID, userID primitive.ObjectID, key string) error {
	f := bson.D{{Key: "_id", Value: storyID}}
	u := bson.D{{Key: key, Value: bson.D{{Key: "likedUserIds", Value: userID}}}}
	_, err := r.Collection.UpdateOne(context.Background(), f, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateCommentID(storyID, commentID primitive.ObjectID, key string) error {
	storyFilter := bson.D{{Key: "_id", Value: storyID}}
	update := bson.D{{Key: key, Value: bson.D{{Key: "commentIds", Value: commentID}}}}
	_, err := r.Collection.UpdateOne(context.Background(), storyFilter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindStories(storyIDs *[]primitive.ObjectID) (*[]model.Story, error) {
	stories := make([]model.Story, 0)
	f := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: storyIDs}}}}
	o := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	c, err := r.Collection.Find(context.Background(), f, o)
	if err != nil {
		return nil, err
	}
	if err = c.All(context.Background(), &stories); err != nil {
		return nil, err
	}
	return &stories, err
}

func (r *repository) FindStoryByID(storyID primitive.ObjectID) (*model.Story, error) {
	story := new(model.Story)
	f := bson.D{{Key: "_id", Value: storyID}}
	singleResult := r.Collection.FindOne(context.Background(), f)
	if err := singleResult.Err(); err != nil {
		return nil, err
	}
	singleResult.Decode(story)
	return story, nil
}

func (r *repository) UpdateViewCount(storyID primitive.ObjectID) (*model.Story, error) {
	story := new(model.Story)
	f := bson.D{{Key: "_id", Value: storyID}}
	u := bson.D{{Key: "$inc", Value: bson.D{{Key: "viewCount", Value: 1}}}}
	storyResult := r.Collection.FindOneAndUpdate(context.Background(), f, u)
	if err := storyResult.Err(); err != nil {
		return nil, err
	}
	storyResult.Decode(story)
	return story, nil
}

func (r *repository) FindRecentStories(timestamp int64) (*[]model.Story, error) {
	stories := make([]model.Story, 0)
	o := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(20)
	f := bson.D{{Key: "isPublished", Value: true}, {Key: "editorsPick", Value: false}, {Key: "createdAt", Value: bson.D{{Key: "$lt", Value: timestamp}}}}
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
