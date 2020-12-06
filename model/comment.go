package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Comment model
type Comment struct {
	ID        string             `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt int64              `json:"createdAt" bson:"createdAt"`
	StoryID   primitive.ObjectID `json:"storyId" bson:"storyId"`
	CreatorID primitive.ObjectID `json:"creatorId" bson:"creatorId"`
	Text      string             `json:text`
}
