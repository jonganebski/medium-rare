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
	Level          int8   `json:"level,omitempty"`
	Text           string `json:"text,omitempty"`
	Code           string `json:"code,omitempty"`
	Caption        string `json:"caption,omitempty"`
	File           file   `json:"file,omitempty"`
	Stretched      bool   `json:"stretched,omitempty"`
	WithBackground bool   `json:"withBackground,omitempty"`
	WithBorder     bool   `json:"withBorder,omitempty"`
}

type file struct {
	URL string `json:"url"`
}
