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

// AuthRouter holds urls for authentication
func AuthRouter(app fiber.Router, userService user.Service) {
	app.Post("/signup", signup(userService))
	app.Post("/signin", signin(userService))
	app.Get("/signout", middleware.Protected, signout())
}

func signup(userService user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := new(model.User)
		email := c.FormValue("email")
		password := c.FormValue("password")
		user.Email = email
		// --- check same email exists or not ---
		_, err := userService.FindUserByEmail(user)
		if err == nil {
			return c.Status(400).SendString("This email already exists")
		}
		// --- make basic user struct ---
		user.ID = "" // this will let mongodb make ID automatically
		user.Password = util.HashPassword(password)
		user.Username = strings.Split(email, "@")[0] // first username is gonna be from email. username is not unique field.
		user.AvatarURL = "https://medium-rare.s3.amazonaws.com/blank-profile.webp"
		user.CreatedAt = time.Now().Unix()
		user.UpdatedAt = time.Now().Unix()
		user.CommentIDs = &[]primitive.ObjectID{} // without these mongodb makes these slice fields null
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
		// --- insert user into mongodb ---
		userOID, err := userService.CreateUser(user)
		if err != nil {
			fmt.Println(err)
			return c.Status(500).SendString("Sorry.. server has a problem")
		}
		// --- generate and issue cookie ---
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
		// --- find signing in user by email ---
		foundUser, err := userService.FindUserByEmail(user)
		if err != nil {
			return c.SendStatus(404)
		}
		// --- verify password ---
		isValid := util.VerifyPassword(password, foundUser.Password)
		if !isValid {
			return c.SendStatus(400)
		}
		// --- generate and issue cookie
		exp := time.Hour * 24 * 7 // 7 days
		userOID, err := primitive.ObjectIDFromHex(foundUser.ID)
		if err != nil {
			return c.SendStatus(500)
		}
		cookie, err := util.GenerateCookieBeta(&userOID, exp)
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
