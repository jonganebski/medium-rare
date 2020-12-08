package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// User model
type User struct {
	ID            string                `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt     int64                 `json:"createdAt"`
	UpdatedAt     int64                 `json:"updatedAt"`
	AvatarURL     string                `json:"avatarUrl,omitempty" bson:"avatarUrl,omitempty"`
	Username      string                `json:"username" bson:"username"`
	Email         string                `json:"email" json:"email"`
	Password      string                `json:"password"`
	Bio           string                `json:"bio,omitempty"`
	FollowerIDs   *[]primitive.ObjectID `json:"followerIds" bson:"followerIds"`
	FollowingIDs  *[]primitive.ObjectID `json:"followingIds" bson:"followingIds"`
	CommentIDs    *[]primitive.ObjectID `json:"commentIds" bson:"commentIds"`
	StoryIDs      *[]primitive.ObjectID `json:"storyIds" bson:"storyIds"`
	LikedStoryIDs *[]primitive.ObjectID `json:"likedStoryIds" bson:"likedStoryIds"`
	SavedStoryIDs *[]primitive.ObjectID `json:"savedStoryIds" bson:"savedStoryIds"`
}
