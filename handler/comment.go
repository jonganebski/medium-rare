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

type commentOutput struct {
	Username     string `json:"username"`
	AvatarURL    string `json:"avatarUrl"`
	CreatedAt    int64  `json:"createdAt"`
	Text         string `json:"text"`
	IsAuthorized bool   `json:"isAuthorized"`
}

// AddComment creates a new comment
func AddComment(c *fiber.Ctx) error {

	userCollection := mg.Db.Collection(UserCollection)
	storyCollection := mg.Db.Collection(StoryCollection)
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

	// --- find added comment ---

	foundComment := new(model.Comment)
	commentFilter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
	commentResult := commentCollection.FindOne(c.Context(), commentFilter)
	if commentResult.Err() != nil {
		return c.SendStatus(404)
	}
	commentResult.Decode(foundComment)

	// --- add commentID into story's commentIDs

	foundCommentOID, err := primitive.ObjectIDFromHex(foundComment.ID)
	if err != nil {
		return c.SendStatus(500)
	}
	storyFilter := bson.D{{Key: "_id", Value: storyOID}}
	storyUpdate := bson.D{{Key: "$push", Value: bson.D{{Key: "commentIds", Value: foundCommentOID}}}}
	storyUpdateResult := storyCollection.FindOneAndUpdate(c.Context(), storyFilter, storyUpdate)
	if storyUpdateResult.Err() != nil {
		return c.SendStatus(500)
	}

	// --- build output ---

	output := new(commentOutput)
	output.AvatarURL = user.AvatarURL
	output.CreatedAt = foundComment.CreatedAt
	output.Text = foundComment.Text
	output.Username = user.Username
	output.IsAuthorized = true

	return c.Status(201).JSON(output)
}

// ProvideComments gets comments of the story and returns them in certain output type
func ProvideComments(c *fiber.Ctx) error {

	storyCollection := mg.Db.Collection(StoryCollection)
	userCollection := mg.Db.Collection(UserCollection)
	commentCollection := mg.Db.Collection(CommentCollection)
	storyID := c.Params("storyId")
	storyOID, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		return c.SendStatus(400)
	}

	var outputItem = new(commentOutput)
	var output []commentOutput = make([]commentOutput, 0)

	// --- commentIDs of the story ---

	story := new(model.Story)
	filter := bson.D{{Key: "_id", Value: storyOID}}
	storyResult := storyCollection.FindOne(c.Context(), filter)
	if storyResult.Err() != nil {
		return c.SendStatus(404)
	}
	storyResult.Decode(story)

	// --- find each comment's creator's username and their avatar id

	comment := new(model.Comment)
	user := new(model.User)

	for _, commentID := range *story.CommentIDs {
		// find comment
		filter = bson.D{{Key: "_id", Value: commentID}}
		commentResult := commentCollection.FindOne(c.Context(), filter)
		if commentResult.Err() != nil {
			return c.SendStatus(404)
		}
		commentResult.Decode(comment)
		// find user
		filter = bson.D{{Key: "_id", Value: comment.CreatorID}}
		userResult := userCollection.FindOne(c.Context(), filter)
		if userResult.Err() != nil {
			return c.SendStatus(404)
		}
		userResult.Decode(user)
		// append to output
		outputItem.Username = user.Username
		outputItem.AvatarURL = user.AvatarURL
		outputItem.CreatedAt = comment.CreatedAt
		outputItem.Text = comment.Text
		outputItem.IsAuthorized = (user.ID == c.Locals("userId"))
		output = append(output, *outputItem)
	}

	return c.Status(200).JSON(output)
}
