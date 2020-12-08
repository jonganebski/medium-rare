package handler

import (
	"fmt"
	"home/jonganebski/github/medium-rare/config"
	"home/jonganebski/github/medium-rare/database"
	"home/jonganebski/github/medium-rare/model"
	"home/jonganebski/github/medium-rare/util"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mg = &database.Mongo

// UserCollection is users collection name
var UserCollection = config.Config("COLLECTION_USER")

// CreateUser creates a user
func CreateUser(c *fiber.Ctx) error {

	userCollection := mg.Db.Collection(UserCollection)
	email := c.FormValue("email")
	password := c.FormValue("password")

	filter := bson.D{{Key: "email", Value: email}}
	userResult := userCollection.FindOne(c.Context(), filter)
	if err := userResult.Err(); err == nil {
		return c.Status(400).SendString("This email already exists")
	}

	user := new(model.User)
	user.ID = ""
	user.Email = email
	user.Password = util.HashPassword(password)
	user.Username = strings.Split(email, "@")[0]
	user.AvatarURL = "http://localhost:4000/image/blank-profile.webp"
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()
	user.CommentIDs = &[]primitive.ObjectID{}
	user.FollowerIDs = &[]primitive.ObjectID{}
	user.FollowingIDs = &[]primitive.ObjectID{}
	user.StoryIDs = &[]primitive.ObjectID{}
	user.LikedStoryIDs = &[]primitive.ObjectID{}
	user.SavedStoryIDs = &[]primitive.ObjectID{}
	user.IsEditor = false

	editorEmails := strings.Fields(config.Config("EDITORS"))
	for _, editorEmail := range editorEmails {
		if editorEmail == email {
			user.IsEditor = true
		}
	}

	insertionResult, err := userCollection.InsertOne(c.Context(), user)
	if err != nil {
		return c.Status(500).SendString("Sorry.. server has a problem")
	}
	filter = bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
	createRecord := userCollection.FindOne(c.Context(), filter)

	createdUser := new(model.User)
	createRecord.Decode(createdUser)

	exp := time.Hour * 24 * 7 // 7 days

	cookie, err := util.GenerateCookie(createdUser, exp)
	if err != nil {
		return c.Status(500).SendString("Sorry.. server has a problem")
	}

	c.Cookie(cookie)

	return c.SendStatus(201)
}

// Signin verify user password and gives jwt token
func Signin(c *fiber.Ctx) error {
	fmt.Println(c.Path())
	userCollection := mg.Db.Collection(UserCollection)
	email := c.FormValue("email")
	password := c.FormValue("password")

	user := new(model.User)
	filter := bson.D{{Key: "email", Value: email}}
	userResult := userCollection.FindOne(c.Context(), filter)
	if err := userResult.Err(); err != nil {
		return c.SendStatus(400)
	}

	userResult.Decode(user)

	isValid := util.VerifyPassword(password, user.Password)
	if !isValid {
		return c.SendStatus(400)
	}

	exp := time.Hour * 24 * 7 // 7 days

	cookie, err := util.GenerateCookie(user, exp)
	if err != nil {
		return c.SendStatus(500)
	}

	c.Cookie(cookie)

	return c.SendStatus(200)
}

// Signout destories cookie
func Signout(c *fiber.Ctx) error {
	c.ClearCookie()
	return c.Redirect("/")
}

// Follow adds author's userID into current user's FollowingIDs and adds currnet userID into author's FollowerIDs
func Follow(c *fiber.Ctx) error {

	userCollection := mg.Db.Collection(UserCollection)

	authorID := c.Params("authorId")
	authorOID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		return c.SendStatus(500)
	}
	userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
	if err != nil {
		return c.SendStatus(500)
	}

	// --- add author's userID into current user's FollowingIDs ---

	filter := bson.D{{Key: "_id", Value: userOID}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "followingIds", Value: authorOID}}}}
	updateResult := userCollection.FindOneAndUpdate(c.Context(), filter, update)
	if updateResult.Err() != nil {
		return c.SendStatus(404)
	}

	// --- add current user's userID into author's FollowingIDs ---

	filter = bson.D{{Key: "_id", Value: authorOID}}
	update = bson.D{{Key: "$push", Value: bson.D{{Key: "followerIds", Value: userOID}}}}
	updateResult = userCollection.FindOneAndUpdate(c.Context(), filter, update)
	if updateResult.Err() != nil {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)
}

// Unfollow unfollow
func Unfollow(c *fiber.Ctx) error {

	userCollection := mg.Db.Collection(UserCollection)

	authorID := c.Params("authorId")
	authorOID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		return c.SendStatus(500)
	}
	userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
	if err != nil {
		return c.SendStatus(500)
	}

	// --- add author's userID into current user's FollowingIDs ---

	filter := bson.D{{Key: "_id", Value: userOID}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "followingIds", Value: authorOID}}}}
	updateResult := userCollection.FindOneAndUpdate(c.Context(), filter, update)
	if updateResult.Err() != nil {
		return c.SendStatus(404)
	}

	// --- add current user's userID into author's FollowingIDs ---

	filter = bson.D{{Key: "_id", Value: authorOID}}
	update = bson.D{{Key: "$pull", Value: bson.D{{Key: "followerIds", Value: userOID}}}}
	updateResult = userCollection.FindOneAndUpdate(c.Context(), filter, update)
	if updateResult.Err() != nil {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)
}

// SettingsPage renders settings page
func SettingsPage(c *fiber.Ctx) error {

	userCollection := mg.Db.Collection(UserCollection)

	userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(500)
	}

	user := new(model.User)
	filter := bson.D{{Key: "_id", Value: userOID}}
	singleResult := userCollection.FindOne(c.Context(), filter)
	if singleResult.Err() != nil {
		return c.SendStatus(404)
	}
	singleResult.Decode(user)

	return c.Render("settings", fiber.Map{"path": c.Path(), "userId": c.Locals("userId"), "userAvatarUrl": user.AvatarURL, "username": user.Username, "userEmail": user.Email, "bio": user.Bio}, "layout/main")
}

// EditUsername updates user's username
func EditUsername(c *fiber.Ctx) error {

	type editUsernameInput struct {
		Username string `json:"username"`
	}

	userCollection := mg.Db.Collection(UserCollection)

	userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
	if err != nil {
		return c.SendStatus(500)
	}

	input := new(editUsernameInput)
	if err := c.BodyParser(input); err != nil {
		return c.SendStatus(400)
	}

	user := new(model.User)
	filter := bson.D{{Key: "_id", Value: userOID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "username", Value: input.Username}}}}
	_, err = userCollection.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return c.SendStatus(500)
	}

	singleResult := userCollection.FindOne(c.Context(), filter)
	if singleResult.Err() != nil {
		return c.SendStatus(404)
	}

	singleResult.Decode(user)

	return c.Status(200).SendString(user.Username)
}

// EditBio updates user's username
func EditBio(c *fiber.Ctx) error {

	type editBioInput struct {
		Bio string `json:"bio"`
	}

	userCollection := mg.Db.Collection(UserCollection)

	userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
	if err != nil {
		return c.SendStatus(500)
	}

	input := new(editBioInput)
	if err := c.BodyParser(input); err != nil {
		return c.SendStatus(400)
	}

	user := new(model.User)
	filter := bson.D{{Key: "_id", Value: userOID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "bio", Value: input.Bio}}}}
	_, err = userCollection.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return c.SendStatus(500)
	}

	singleResult := userCollection.FindOne(c.Context(), filter)
	if singleResult.Err() != nil {
		return c.SendStatus(404)
	}

	singleResult.Decode(user)

	return c.Status(200).SendString(user.Bio)
}

// EditUserAvatar updates user's username
func EditUserAvatar(c *fiber.Ctx) error {

	type editUserAvatarInput struct {
		AvatarURL string `json:"avatarUrl"`
	}

	userCollection := mg.Db.Collection(UserCollection)

	userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
	if err != nil {
		return c.SendStatus(500)
	}

	input := new(editUserAvatarInput)
	if err := c.BodyParser(input); err != nil {
		return c.SendStatus(400)
	}

	user := new(model.User)
	filter := bson.D{{Key: "_id", Value: userOID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "avatarUrl", Value: input.AvatarURL}}}}
	_, err = userCollection.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return c.SendStatus(500)
	}

	singleResult := userCollection.FindOne(c.Context(), filter)
	if singleResult.Err() != nil {
		return c.SendStatus(404)
	}

	singleResult.Decode(user)

	return c.Status(200).SendString(user.AvatarURL)
}

// DeleteUser removes current user and all related documents from the database and related fields
func DeleteUser(c *fiber.Ctx) error {

	// --- check password ---

	// --- delete user's stories ---

	// --- delete user's comments ---

	// --- delte from other users's followingIDs ---

	// --- delete avatar photo ---

	// --- delete user ---

	return c.Redirect("/")
}
