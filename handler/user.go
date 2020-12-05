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
		fmt.Println("this email already exists")
		return c.SendStatus(400)
	}

	user := new(model.User)
	user.ID = ""
	user.Email = email
	user.Password = util.HashPassword(password)
	user.Username = strings.Split(email, "@")[0]
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()
	user.CommentIDs = &[]primitive.ObjectID{}
	user.FollowerIDs = &[]primitive.ObjectID{}
	user.FollowingIDs = &[]primitive.ObjectID{}
	user.StoryIDs = &[]primitive.ObjectID{}
	user.LikedStoryIDs = &[]primitive.ObjectID{}

	insertionResult, err := userCollection.InsertOne(c.Context(), user)
	if err != nil {
		fmt.Println("Failed to add to database")
		c.SendStatus(500)
	}
	filter = bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
	createRecord := userCollection.FindOne(c.Context(), filter)

	createdUser := new(model.User)
	createRecord.Decode(createdUser)

	exp := time.Hour * 24 * 7 // 7 days

	cookie, err := util.GenerateCookie(createdUser, exp)
	if err != nil {
		fmt.Println("error at generating cookie")
		return c.SendStatus(500)
	}

	c.Cookie(cookie)

	return c.SendStatus(201)
}

// Signin verify user password and gives jwt token
func Signin(c *fiber.Ctx) error {

	userCollection := mg.Db.Collection(UserCollection)
	email := c.FormValue("email")
	password := c.FormValue("password")

	user := new(model.User)
	filter := bson.D{{Key: "email", Value: email}}
	userResult := userCollection.FindOne(c.Context(), filter)
	if err := userResult.Err(); err != nil {
		fmt.Println("failed to find email")
		return c.SendStatus(400)
	}

	userResult.Decode(user)

	isValid := util.VerifyPassword(password, user.Password)
	if !isValid {
		fmt.Println("wrong password")
		return c.SendStatus(400)
	}

	exp := time.Hour * 24 * 7 // 7 days

	cookie, err := util.GenerateCookie(user, exp)
	if err != nil {
		fmt.Println("error at generating cookie")
		return c.SendStatus(500)
	}

	c.Cookie(cookie)

	return c.SendStatus(200)
}
