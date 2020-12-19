package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Story model
type Story struct {
	ID           string                `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt    int64                 `json:"createdAt" bson:"createdAt"`
	UpdatedAt    int64                 `json:"updatedAt" bson:"updatedAt"`
	CreatorID    primitive.ObjectID    `json:"creatorId" bson:"creatorId"`
	Blocks       []Block               `json:"blocks"`
	EditorjsVer  string                `json:"version"`
	ViewCount    uint32                `json:"viewCount" bson:"viewCount"`
	EditorsPick  bool                  `json:"editorsPick" bson:"editorsPick"`
	IsPublished  bool                  `json:"isPublished" bson:"isPublished"`
	LikedUserIDs *[]primitive.ObjectID `json:"likedUserIds" bson:"likedUserIds"`
	CommentIDs   *[]primitive.ObjectID `json:"commentIds" bson:"commentIds"`
}

// Block struct
type Block struct {
	Type string `json:"type"`
	Data data   `json:"data"`
}

type data struct {
	Alignment      string   `json:"alignment,omiempty" bson:"alignment,omiempty"`
	Level          int8     `json:"level,omitempty" bson:"level,omitempty"`
	Text           string   `json:"text,omitempty" bson:"text,omitempty"`
	Code           string   `json:"code,omitempty" bson:"code,omitempty"`
	Caption        string   `json:"caption,omitempty" bson:"caption,omitempty"`
	File           file     `json:"file,omitempty" bson:"file,omitempty"`
	Stretched      bool     `json:"stretched,omitempty" bson:"stretched,omitempty"`
	Style          string   `json:"style,omitempty" bson:"style,omitempty"`
	Items          []string `json:"items,omitempty" bson:"items,omitempty"`
	WithBackground bool     `json:"withBackground,omitempty" bson:"withBackground,omitempty"`
	WithBorder     bool     `json:"withBorder,omitempty" bson:"withBorder,omitempty"`
}

type file struct {
	URL string `json:"url,omitempty" bson:"url,omitempty"`
}

// StoryCardOutput describes storyCard.pug partial
type StoryCardOutput struct {
	StoryID        string `json:"storyId"`
	AuthorID       string `json:"authorId"`
	AuthorUsername string `json:"authorUsername"`
	CreatedAt      int64  `json:"createdAt"`
	UpdatedAt      int64  `json:"updatedAt"`
	Header         string `json:"header"`
	Body           string `json:"body"`
	CoverImgURL    string `json:"coverImgUrl"`
	ReadTime       string `json:"readTime"`
	Ranking        int    `json:"ranking,omitempty"`
	IsPublished    bool   `json:"isPublished"`
}
