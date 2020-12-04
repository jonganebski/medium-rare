package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Story model
type Story struct {
	ID           string                `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt    int64                 `json:"createdAt"`
	UpdatedAt    int64                 `json:"updatedAt"`
	CreatorID    primitive.ObjectID    `json:"creatorId" bson:"creatorId"`
	Text         string                `json:"text"`
	ViewCount    uint32                `json:"viewCount" bson:"viewCount"`
	LikedUserIDs *[]primitive.ObjectID `json:"likedUserIds" bson:"likedUserIds"`
	CommentIDs   *[]primitive.ObjectID `json:"commentIds" bson:"commentIds"`
}
