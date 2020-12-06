package handler

import (
	"fmt"
	"home/jonganebski/github/medium-rare/config"
	"home/jonganebski/github/medium-rare/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CommentCollection is comments collection name
var CommentCollection = config.Config("COLLECTION_COMMENT")

// AddComment creates a new comment
func AddComment(c *fiber.Ctx) error {

	type addCommentOutput struct {
		Username  string `json:"username"`
		CreatedAt int64  `json:"createdAt"`
		AvatarURL string `json:"avatarUrl"`
		Text      string `json:"text"`
	}

	userCollection := mg.Db.Collection(UserCollection)
	commentCollection := mg.Db.Collection(CommentCollection)

	comment := new(model.Comment)
	if err := c.BodyParser(comment); err != nil {
		return c.SendStatus(400)
	}
	storyID := c.Params("storyId")

	userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
	if err != nil {
		return c.SendStatus(500)
	}
	storyOID, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		return c.SendStatus(500)
	}

	// --- find user ---

	user := new(model.User)
	userFilter := bson.D{{Key: "_id", Value: userOID}}
	userResult := userCollection.FindOne(c.Context(), userFilter)
	if userResult.Err() != nil {
		return c.SendStatus(404)
	}
	userResult.Decode(user)

	// --- build comment struct properly

	comment.ID = ""
	comment.CreatorID = userOID
	comment.CreatedAt = time.Now().Unix()
	comment.StoryID = storyOID

	// --- add comment in mongoDB ---

	insertionResult, err := commentCollection.InsertOne(c.Context(), comment)
	if err != nil {
		return c.SendStatus(500)
	}

	// find added comment

	foundComment := new(model.Comment)
	commentFilter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
	commentResult := commentCollection.FindOne(c.Context(), commentFilter)
	if commentResult.Err() != nil {
		return c.SendStatus(404)
	}
	commentResult.Decode(foundComment)

	// --- build output

	output := new(addCommentOutput)
	output.AvatarURL = user.AvatarURL
	output.CreatedAt = foundComment.CreatedAt
	output.Text = foundComment.Text
	output.Username = user.Username

	return c.Status(201).JSON(output)
}
