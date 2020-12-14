package routes

import (
	"fmt"
	"home/jonganebski/github/medium-rare/package/story"
	"home/jonganebski/github/medium-rare/package/user"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdminRouter has api routes only for editors
func AdminRouter(admin fiber.Router, userService user.Service, storyService story.Service) {
	admin.Patch("/pick/:storyId", pickStory(userService, storyService))
	admin.Patch("/unpick/:storyId", unpickStory(userService, storyService))
}

func unpickStory(userService user.Service, storyService story.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		storyID := c.Params("storyId")
		storyOID, err := primitive.ObjectIDFromHex(storyID)
		if err != nil {
			return c.SendStatus(500)
		}
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}

		// --- find current user and check he/she is editor ---
		currentUser, err := userService.FindUserByID(userOID)
		if !currentUser.IsEditor {
			return c.SendStatus(403)
		}

		// --- update story's EditorsPick field false ---
		err = storyService.UnpickStory(storyOID)
		if err != nil {
			return c.SendStatus(500)
		}

		return c.SendStatus(200)
	}
}

func pickStory(userService user.Service, storyService story.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		storyID := c.Params("storyId")
		storyOID, err := primitive.ObjectIDFromHex(storyID)
		if err != nil {
			return c.SendStatus(500)
		}
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}

		// --- find current user and check he/she is editor ---
		currentUser, err := userService.FindUserByID(userOID)
		if !currentUser.IsEditor {
			return c.SendStatus(403)
		}

		// --- update story's EditorsPick field true ---
		err = storyService.PickStory(storyOID)
		if err != nil {
			return c.SendStatus(500)
		}

		return c.SendStatus(200)
	}
}
