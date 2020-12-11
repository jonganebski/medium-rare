package routes

import (
	"fmt"
	"home/jonganebski/github/medium-rare/middleware"
	"home/jonganebski/github/medium-rare/model"
	"home/jonganebski/github/medium-rare/package/comment"
	"home/jonganebski/github/medium-rare/package/story"
	"home/jonganebski/github/medium-rare/package/user"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type commentOutput struct {
	CommentID    string `json:"commentId"`
	Username     string `json:"username"`
	AvatarURL    string `json:"avatarUrl"`
	CreatedAt    int64  `json:"createdAt"`
	Text         string `json:"text"`
	IsAuthorized bool   `json:"isAuthorized"`
}

// CommentRouter has api routes for comment
func CommentRouter(app fiber.Router, userService user.Service, storyService story.Service, commentService comment.Service) {
	app.Get("/comment/:storyId", provideComments(userService, storyService, commentService))
	app.Post("/comment/:storyId", middleware.APIGuard, addComment(userService, storyService, commentService))
}

func addComment(userService user.Service, storyService story.Service, commentService comment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		comment := new(model.Comment)
		if err := c.BodyParser(comment); err != nil {
			return c.SendStatus(400)
		}
		storyID := c.Params("storyId")

		storyOID, err := primitive.ObjectIDFromHex(storyID)
		if err != nil {
			return c.SendStatus(500)
		}
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}

		currentUser, err := userService.FindUserByID(userOID)
		if err != nil {
			return c.Status(404).SendString("User not found")
		}

		comment.ID = ""
		comment.CreatorID = userOID
		comment.CreatedAt = time.Now().Unix()
		comment.StoryID = storyOID

		comment, err = commentService.CreateComment(comment)
		if err != nil {
			return c.Status(500).SendString("Create comment failed")
		}
		commentOID, err := primitive.ObjectIDFromHex(comment.ID)
		if err != nil {
			return c.SendStatus(500)
		}
		err = userService.AddCommentID(userOID, commentOID)
		if err != nil {
			return c.Status(500).SendString("Update failed")
		}
		err = storyService.AddCommentID(storyOID, commentOID)
		if err != nil {
			return c.Status(500).SendString("Update failed")
		}

		output := new(commentOutput)
		output.CommentID = comment.ID
		output.AvatarURL = currentUser.AvatarURL
		output.CreatedAt = comment.CreatedAt
		output.Text = comment.Text
		output.Username = currentUser.Username
		output.IsAuthorized = true

		return c.Status(201).JSON(output)
	}
}

func provideComments(userService user.Service, storyService story.Service, commentService comment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		storyID := c.Params("storyId")
		storyOID, err := primitive.ObjectIDFromHex(storyID)
		if err != nil {
			return c.SendStatus(400)
		}
		story, err := storyService.FindStoryByID(storyOID)
		if err != nil {
			return c.Status(404).SendString("Story not found")
		}
		comments, err := commentService.FindComments(story.CommentIDs)
		if err != nil {
			return c.Status(404).SendString("Comments not found")
		}

		outputItem := new(commentOutput)
		output := make([]commentOutput, 0)

		for _, comment := range *comments {
			creator, err := userService.FindUserByID(comment.CreatorID)
			if err != nil {
				return c.Status(404).SendString("Commenter not found")
			}
			outputItem.CommentID = comment.ID
			outputItem.Username = creator.Username
			outputItem.AvatarURL = creator.AvatarURL
			outputItem.CreatedAt = comment.CreatedAt
			outputItem.Text = comment.Text
			outputItem.IsAuthorized = (creator.ID == c.Locals("userId"))
			output = append(output, *outputItem)
		}

		return c.Status(200).JSON(output)
	}
}
