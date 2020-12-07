package handler

import (
	"fmt"
	"home/jonganebski/github/medium-rare/config"
	"home/jonganebski/github/medium-rare/model"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StoryCollection is stories collection name
var StoryCollection = config.Config("COLLECTION_STORY")

// Home renders homepage
func Home(c *fiber.Ctx) error {

	type homeOutput struct {
		StoryID        string `json:"storyId"`
		AuthorUsername string `json:"authorUsername"`
		CreatedAt      int64  `json:"createdAt"`
		Header         string `json:"header"`
		Body           string `json:"body"`
		CoverImgURL    string `json:"coverImgUrl"`
	}

	storyCollection := mg.Db.Collection(StoryCollection)
	userCollection := mg.Db.Collection(UserCollection)

	outputItem := new(homeOutput)
	output := make([]homeOutput, 0)

	storyFindOptions := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(20)
	storyFilter := bson.D{{}}
	cursor, err := storyCollection.Find(c.Context(), storyFilter, storyFindOptions)
	if err != nil {
		fmt.Println("Error at finding stories")
		return c.SendStatus(500)
	}

	stories := make([]model.Story, 0)

	if err := cursor.All(c.Context(), &stories); err != nil {
		fmt.Println("error at cursor iteration")
		fmt.Println(err)
		return c.SendStatus(500)
	}

	author := new(model.User)
	for _, story := range stories {

		// --- find body & coverImgUrl---

		body := ""
		coverImgURL := ""
		for _, block := range story.Blocks {
			if block.Type == "paragraph" && body == "" {
				body = block.Data.Text
			}
			if block.Type == "image" && coverImgURL == "" {
				coverImgURL = block.Data.File.URL
			}
			if body != "" && coverImgURL != "" {
				break
			}
		}

		// --- find author ---

		authorFilter := bson.D{{Key: "_id", Value: story.CreatorID}}
		authorResult := userCollection.FindOne(c.Context(), authorFilter)
		if authorResult.Err() != nil {
			fmt.Println("user does not exist")
			return c.SendStatus(500)
		}
		authorResult.Decode(author)

		// --- build outputItem and append to output ---

		outputItem.AuthorUsername = author.Username
		outputItem.StoryID = story.ID
		outputItem.Header = story.Blocks[0].Data.Text
		outputItem.Body = body
		outputItem.CreatedAt = story.CreatedAt
		outputItem.CoverImgURL = coverImgURL
		output = append(output, *outputItem)
	}

	return c.Render("home", fiber.Map{"path": c.Path(), "userId": c.Locals("userId"), "output": output}, "layout/main")
}

// NewStory renders a page where a user writes a new story
func NewStory(c *fiber.Ctx) error {
	return c.Render("newStory", fiber.Map{"path": c.Path(), "userId": c.Locals("userId")}, "layout/main")
}

// ReadStory renders a page where a user reads a story
func ReadStory(c *fiber.Ctx) error {
	// status 400 404 구분할 것
	storyCollection := mg.Db.Collection(StoryCollection)
	userCollection := mg.Db.Collection(UserCollection)

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

	// --- currnet user ---
	userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
	if err != nil {
		return c.SendStatus(400)
	}
	filter = bson.D{{Key: "_id", Value: userOID}}
	userResult := userCollection.FindOne(c.Context(), filter)
	if userResult.Err() != nil {
		return c.SendStatus(400)
	}

	user := new(model.User)
	userResult.Decode(user)

	// --- did currnet user liked this story? ---

	didLiked := false
	for _, likedStoryID := range *user.LikedStoryIDs {
		if likedStoryID == storyOID {
			didLiked = true
			break
		}
	}

	// --- did current user bookmarked this story?

	bookmarked := false
	for _, savedStoryID := range *user.SavedStoryIDs {
		if savedStoryID == storyOID {
			bookmarked = true
			break
		}
	}

	// --- is current user following author of the story

	isFollowing := false
	for _, followerID := range *author.FollowerIDs {
		if followerID == userOID {
			isFollowing = true
			break
		}
	}

	fmt.Println(isFollowing)

	return c.Render("readStory",
		fiber.Map{"path": c.Path(),
			"userId":      c.Locals("userId"),
			"username":    user.Username,
			"story":       story,
			"author":      author,
			"didLiked":    didLiked,
			"bookmarked":  bookmarked,
			"isFollowing": isFollowing},
		"layout/main")
}

// EditStory renders a page where a user edits his/her story
func EditStory(c *fiber.Ctx) error {

	storyCollection := mg.Db.Collection(StoryCollection)

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

	userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
	if err != nil {
		fmt.Println("error at conversion")
		return c.SendStatus(500)
	}

	if userOID != story.CreatorID {
		fmt.Println("You are not authorized")
		c.Redirect("/")
	}

	return c.Render("editStory", fiber.Map{"path": c.Path(), "userId": c.Locals("userId"), "story": story}, "layout/main")
}

// ProvideStoryBlocks returns blocks of the story
func ProvideStoryBlocks(c *fiber.Ctx) error {

	storyCollection := mg.Db.Collection(StoryCollection)

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

	storyCollection := mg.Db.Collection(StoryCollection)
	userCollection := mg.Db.Collection(UserCollection)

	// --- story ---
	story := new(model.Story)

	if err := c.BodyParser(story); err != nil {
		return c.SendStatus(400)
	}

	userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
	if err != nil {
		return c.SendStatus(500)
	}

	story.ID = ""
	story.CreatedAt = time.Now().Unix()
	story.UpdatedAt = time.Now().Unix()
	story.CreatorID = userOID
	story.LikedUserIDs = &[]primitive.ObjectID{}
	story.CommentIDs = &[]primitive.ObjectID{}
	story.ViewCount = 0

	insertiionResult, err := storyCollection.InsertOne(c.Context(), story)
	if err != nil {
		return c.SendStatus(500)
	}

	// --- user ---

	filter := bson.D{{Key: "_id", Value: userOID}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "storyIds", Value: insertiionResult.InsertedID}}}}
	updateResult := userCollection.FindOneAndUpdate(c.Context(), filter, update)
	if updateResult.Err() != nil {
		return c.SendStatus(500)
	}

	stringInsertedID := insertiionResult.InsertedID.(primitive.ObjectID).Hex()

	return c.Status(201).SendString(stringInsertedID)
}

// UpdateStory updates story
func UpdateStory(c *fiber.Ctx) error {

	storyCollection := mg.Db.Collection(StoryCollection)
	uploadedStory := new(model.Story)

	if err := c.BodyParser(uploadedStory); err != nil {
		fmt.Println("error at body parser")
		return c.SendStatus(400)
	}

	storyID := c.Params("storyId")

	userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
	if err != nil {
		fmt.Println("error at conversion")
		return c.SendStatus(500)
	}
	storyOID, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		fmt.Println("error at conversion")
		return c.SendStatus(500)
	}

	filter := bson.D{{Key: "_id", Value: storyOID}}
	findResult := storyCollection.FindOne(c.Context(), filter)
	if findResult.Err() != nil {
		fmt.Println("Story not found")
		return c.SendStatus(400)
	}
	foundStory := new(model.Story)
	findResult.Decode(foundStory)

	if userOID != foundStory.CreatorID {
		fmt.Println("You are not authorized")
		return c.SendStatus(400)
	}

	filter = bson.D{{Key: "_id", Value: storyOID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "blocks", Value: uploadedStory.Blocks}, {Key: "updatedAt", Value: time.Now().Unix()}}}}
	updateResult := storyCollection.FindOneAndUpdate(c.Context(), filter, update)
	if updateResult.Err() != nil {
		fmt.Println("Error at update")
		fmt.Println(updateResult.Err())
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

// HandleLikeCount deals with increment/decreament of like count of the story
func HandleLikeCount(c *fiber.Ctx) error {

	storyCollection := mg.Db.Collection(StoryCollection)
	userCollection := mg.Db.Collection(UserCollection)

	storyID := c.Params("storyId")
	p := c.Params("plusMinus") // bigger than zero -> increase like / smaller than zero -> decrease like
	plusMinus, err := strconv.Atoi(p)
	if err != nil || plusMinus == 0 {
		return c.SendStatus(400)
	}

	storyOID, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		return c.SendStatus(500)
	}
	userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
	if err != nil {
		return c.SendStatus(500)
	}

	// --- find user ---
	user := new(model.User)
	filter := bson.D{{Key: "_id", Value: userOID}}
	singleResult := userCollection.FindOne(c.Context(), filter)
	if singleResult.Err() != nil {
		return c.SendStatus(404)
	}
	singleResult.Decode(user)

	// --- check the user is allowed to like or unlike the story ---
	var isAllowed bool
	if plusMinus > 0 {
		isAllowed = true
		for _, likedStoryID := range *user.LikedStoryIDs {
			if likedStoryID == storyOID {
				isAllowed = false
			}
		}
	}
	if plusMinus < 0 {
		isAllowed = false
		for _, likedStoryID := range *user.LikedStoryIDs {
			if likedStoryID == storyOID {
				isAllowed = true
			}
		}
	}
	if !isAllowed {
		return c.SendStatus(400)
	}

	var key string
	if plusMinus > 0 {
		key = "$push"
	}
	if plusMinus < 0 {
		key = "$pull"
	}

	// -- update user's likedStoryIDs ---

	filter = bson.D{{Key: "_id", Value: userOID}}
	update := bson.D{{Key: key, Value: bson.D{{Key: "likedStoryIds", Value: storyOID}}}}
	updateResult := userCollection.FindOneAndUpdate(c.Context(), filter, update)
	if updateResult.Err() != nil {
		return c.SendStatus(500)
	}

	// --- update story's likedUserIDs ---

	filter = bson.D{{Key: "_id", Value: storyOID}}
	update = bson.D{{Key: key, Value: bson.D{{Key: "likedUserIds", Value: userOID}}}}
	updateResult = storyCollection.FindOneAndUpdate(c.Context(), filter, update)
	if updateResult.Err() != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

// BookmarkStory saves storyID into user's SavedStoryIDs
func BookmarkStory(c *fiber.Ctx) error {

	userCollection := mg.Db.Collection(UserCollection)

	storyID := c.Params("storyId")
	storyOID, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		return c.SendStatus(500)
	}
	userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))

	filter := bson.D{{Key: "_id", Value: userOID}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "savedStoryIds", Value: storyOID}}}}
	updateResult := userCollection.FindOneAndUpdate(c.Context(), filter, update)
	if updateResult.Err() != nil {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)
}

// DisBookmarkStory removes storyID into user's SavedStoryIDs
func DisBookmarkStory(c *fiber.Ctx) error {

	userCollection := mg.Db.Collection(UserCollection)

	storyID := c.Params("storyId")
	storyOID, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		return c.SendStatus(500)
	}
	userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))

	filter := bson.D{{Key: "_id", Value: userOID}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "savedStoryIds", Value: storyOID}}}}
	updateResult := userCollection.FindOneAndUpdate(c.Context(), filter, update)
	if updateResult.Err() != nil {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)
}
