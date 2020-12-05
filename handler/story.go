package handler

import (
	"fmt"
	"home/jonganebski/github/fibersteps-server/config"
	"home/jonganebski/github/fibersteps-server/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var storyCollectionName = config.Config("COLLECTION_STORY")

// Home renders homepage
func Home(c *fiber.Ctx) error {

	return c.Render("home", fiber.Map{"path": c.Path(), "userId": c.Locals("userId")}, "layout/main")
}

// NewStory renders a page where a user writes a new story
func NewStory(c *fiber.Ctx) error {
	return c.Render("newStory", fiber.Map{"path": c.Path(), "userId": c.Locals("userId")}, "layout/main")
}

// ReadStory renders a page where a user reads a story
func ReadStory(c *fiber.Ctx) error {

	storyCollection := mg.Db.Collection(storyCollectionName)
	userCollection := mg.Db.Collection(userCollectionName)

	// --- story ---
	storyID := c.Params("storyId")
	storyOID, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		fmt.Println("error at conversion")
		return c.SendStatus(500)
	}
	filter := bson.D{{Key: "_id", Value: storyOID}}
	storyResult := storyCollection.FindOne(c.Context(), filter)
	if storyResult.Err() != nil {
		fmt.Println("Story does not exist")
		return c.SendStatus(400)
	}

	story := new(model.Story)
	storyResult.Decode(story)

	// --- creator of the story ---
	filter = bson.D{{Key: "_id", Value: story.CreatorID}}
	authorResult := userCollection.FindOne(c.Context(), filter)
	if authorResult.Err() != nil {
		fmt.Println("author does not exist")
		return c.SendStatus(400)
	}

	author := new(model.User)
	authorResult.Decode(author)

	return c.Render("readStory", fiber.Map{"path": c.Path(), "userId": c.Locals("userId"), "story": story, "author": author}, "layout/main")
}

// ProvideStoryBlocks returns blocks of the story
func ProvideStoryBlocks(c *fiber.Ctx) error {

	storyCollection := mg.Db.Collection(storyCollectionName)

	storyID := c.Params("storyId")
	storyOID, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		fmt.Println("error at conversion")
		return c.SendStatus(500)
	}
	filter := bson.D{{Key: "_id", Value: storyOID}}
	storyResult := storyCollection.FindOne(c.Context(), filter)
	if storyResult.Err() != nil {
		fmt.Println("Story does not exist")
		return c.SendStatus(400)
	}

	story := new(model.Story)
	storyResult.Decode(story)

	fmt.Println(story)

	return c.JSON(story.Blocks)
}

// AddStory creates a new story
func AddStory(c *fiber.Ctx) error {

	storyCollection := mg.Db.Collection(storyCollectionName)
	story := new(model.Story)

	if err := c.BodyParser(story); err != nil {
		fmt.Println("error at body parser")
		return c.SendStatus(400)
	}

	userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
	if err != nil {
		fmt.Println("error at conversion")
		return c.SendStatus(500)
	}

	story.ID = ""
	story.CreatedAt = time.Now().Unix()
	story.UpdatedAt = time.Now().Unix()
	story.CreatorID = userOID
	story.LikedUserIDs = &[]primitive.ObjectID{}
	story.CommentIDs = &[]primitive.ObjectID{}
	story.ViewCount = 0

	fmt.Println(story)

	_, err = storyCollection.InsertOne(c.Context(), story)
	if err != nil {
		fmt.Println("error at insertion")
		return c.SendStatus(500)
	}

	return c.Status(201).JSON(story)
}
