package routes

import (
	"fmt"
	"home/jonganebski/github/medium-rare/config"
	"home/jonganebski/github/medium-rare/middleware"
	"home/jonganebski/github/medium-rare/model"
	"home/jonganebski/github/medium-rare/package/user"
	"home/jonganebski/github/medium-rare/util"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRouter has routes related with the user
func UserRouter(app fiber.Router, userService user.Service) {
	app.Post("/signup", signup(userService))
	app.Post("/signin", signin(userService))
	app.Post("/signout", middleware.Protected, signout())
	// route 재설정 필요
	app.Post("/bookmark/:storyId", middleware.APIGuard, bookmarkStory(userService))
	app.Post("/follow/:authorId", middleware.APIGuard, follow(userService))
	app.Post("/unfollow/:authorId", middleware.APIGuard, unfollow(userService))
}

func unfollow(userService user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorID := c.Params("authorId")
		authorOID, err := primitive.ObjectIDFromHex(authorID)
		if err != nil {
			return c.SendStatus(500)
		}
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}
		err = userService.RemoveFollowerID(authorOID, userOID)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}
		err = userService.RemoveFollowingID(userOID, authorOID)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}

		return c.SendStatus(200)
	}
}

func follow(userService user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorID := c.Params("authorId")
		authorOID, err := primitive.ObjectIDFromHex(authorID)
		if err != nil {
			return c.SendStatus(500)
		}
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}
		author, err := userService.FindUserByID(authorOID)
		if err != nil {
			return c.Status(404).SendString("Author not found")
		}

		for _, followerID := range *author.FollowerIDs {
			if followerID == userOID {
				return c.Status(400).SendString("You are already following this user.")
			}
		}

		err = userService.AddFollowingID(userOID, authorOID)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}
		err = userService.AddFollowerID(authorOID, userOID)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}
		return c.SendStatus(200)
	}
}

func bookmarkStory(userService user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		storyID := c.Params("storyId")
		storyOID, err := primitive.ObjectIDFromHex(storyID)
		if err != nil {
			return c.SendStatus(500)
		}
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		err = userService.BookmarkStory(userOID, storyOID)
		if err != nil {
			return err
		}

		return c.SendStatus(200)
	}
}

func signup(userService user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := new(model.User)
		email := c.FormValue("email")
		password := c.FormValue("password")
		user.Email = email
		_, err := userService.FindUserByEmail(user)
		if err == nil {
			return c.Status(400).SendString("This email already exists")
		}
		user.ID = ""
		user.Password = util.HashPassword(password)
		user.Username = strings.Split(email, "@")[0]
		// user.AvatarURL = "http://localhost:4000/image/blank-profile.webp"
		user.AvatarURL = "https://medium-rare.s3.amazonaws.com/blank-profile.webp"
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
		userOID, err := userService.CreateUser(user)
		if err != nil {
			return c.Status(500).SendString("Sorry.. server has a problem")
		}
		exp := time.Hour * 24 * 7 // 7 days
		cookie, err := util.GenerateCookieBeta(userOID, exp)
		if err != nil {
			return c.Status(500).SendString("Sorry.. server has a problem")
		}

		c.Cookie(cookie)

		return c.SendStatus(201)
	}
}

func signin(userService user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := new(model.User)
		email := c.FormValue("email")
		password := c.FormValue("password")
		user.Email = email
		foundUser, err := userService.FindUserByEmail(user)
		if err != nil {
			return c.SendStatus(404)
		}
		isValid := util.VerifyPassword(password, foundUser.Password)
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
}

func signout() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.ClearCookie()
		return c.Redirect("/")
	}
}
