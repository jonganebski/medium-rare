package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// User model
type User struct {
	ID           string                `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt    int64                 `json:"createdAt"`
	UpdatedAt    int64                 `json:"updatedAt"`
	AvatarURL    string                `json:"avatarUrl,omitempty"`
	Firstname    string                `json:"firstname"`
	Lastname     string                `json:"lastname"`
	Username     string                `json:"username" bson:"username"`
	Email        string                `json:"email" json:"email"`
	Password     string                `json:"password"`
	City         string                `json:"city"`
	TimeZone     string                `json:"timeZone" bson:"timeZone"`
	About        string                `json:"about,omitempty"`
	FollowerIDs  *[]primitive.ObjectID `json:"followerIds" bson:"followerIds"`
	FollowingIDs *[]primitive.ObjectID `json:"followingIds" bson:"followingIds"`
	CommentIDs   *[]primitive.ObjectID `json:"commentIds" bson:"commentIds"`
	TripIDs      *[]primitive.ObjectID `json:"tripIds" bson:"tripIds"`
	LikedStepIDs *[]primitive.ObjectID `json:"likedStepIds"`
}
