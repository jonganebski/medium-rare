package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// User model
type User struct {
	ID            string                `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt     int64                 `json:"createdAt" bson:"createdAt"`
	UpdatedAt     int64                 `json:"updatedAt" bson:"updatedAt"`
	AvatarURL     string                `json:"avatarUrl" bson:"avatarUrl"`
	Username      string                `json:"username" bson:"username"`
	Email         string                `json:"email"`
	Password      string                `json:"password"`
	Bio           string                `json:"bio,omitempty"`
	IsEditor      bool                  `json:"isEditor" bson:"isEditor"`
	FollowerIDs   *[]primitive.ObjectID `json:"followerIds,omitempty" bson:"followerIds"`
	FollowingIDs  *[]primitive.ObjectID `json:"followingIds,omitempty" bson:"followingIds"`
	CommentIDs    *[]primitive.ObjectID `json:"commentIds,omitempty" bson:"commentIds"`
	StoryIDs      *[]primitive.ObjectID `json:"storyIds,omitempty" bson:"storyIds"`
	LikedStoryIDs *[]primitive.ObjectID `json:"likedStoryIds,omitempty" bson:"likedStoryIds"`
	SavedStoryIDs *[]primitive.ObjectID `json:"savedStoryIds,omitempty" bson:"savedStoryIds"`
}
