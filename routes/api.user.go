package routes

import (
	"bytes"
	"fmt"
	myaws "home/jonganebski/github/medium-rare/aws"
	"home/jonganebski/github/medium-rare/config"
	"home/jonganebski/github/medium-rare/middleware"
	"home/jonganebski/github/medium-rare/model"
	"home/jonganebski/github/medium-rare/package/comment"
	"home/jonganebski/github/medium-rare/package/story"
	"home/jonganebski/github/medium-rare/package/user"
	"home/jonganebski/github/medium-rare/util"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRouter has routes related with the user
func UserRouter(app fiber.Router, userService user.Service, storyService story.Service, commentService comment.Service) {
	app.Post("/signup", signup(userService))
	app.Post("/signin", signin(userService))
	app.Post("/signout", middleware.Protected, signout())
	// route 재설정 필요
	// bookmark, disbookmark가 post와 delete여야할 이유가 없음.
	app.Post("/bookmark/:storyId", middleware.APIGuard, bookmarkStory(userService))
	app.Post("/follow/:authorId", middleware.APIGuard, follow(userService))
	app.Post("/unfollow/:authorId", middleware.APIGuard, unfollow(userService))
	app.Patch("/user/username", middleware.APIGuard, editUsername(userService))
	app.Patch("/user/bio", middleware.APIGuard, editBio(userService))
	app.Patch("/user/avatar", middleware.APIGuard, editAvatar(userService))
	app.Patch("/user/password", middleware.APIGuard, editPassword(userService))
	app.Delete("/bookmark/:storyId", middleware.APIGuard, disbookmarkStory(userService))
	app.Delete("/user", middleware.APIGuard, removeAccount(userService, storyService, commentService))
}

func removeAccount(userService user.Service, storyService story.Service, commentService comment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type removeAccountInput struct {
			Password string `json:"password"`
		}

		input := new(removeAccountInput)
		if err := c.BodyParser(input); err != nil {
			return c.SendStatus(400)
		}

		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}

		// --- find user ---

		user, err := userService.FindUserByID(userOID)
		if err != nil {
			return c.Status(404).SendString("User not found")
		}

		// --- check password ---

		isValid := util.VerifyPassword(input.Password, user.Password)
		if !isValid {
			return c.Status(400).SendString("You are not authorized.")
		}

		// --- delete user's stories and related fields ---

		for _, storyOID := range *user.StoryIDs {
			story, err := storyService.FindStoryByID(storyOID)
			if err != nil {
				return c.Status(404).SendString("Story not found")
			}
			if story.CreatorID != userOID {
				return c.Status(400).SendString("You are not authorized.")
			}
			err = userService.RemoveManyLikedStoryIDs(storyOID)
			if err != nil {
				return c.Status(500).SendString("Failed to update")
			}
			err = userService.RemoveManySavedStoryIDs(storyOID)
			if err != nil {
				return c.Status(500).SendString("Failed to update")
			}
			err = commentService.RemoveComments(story.CommentIDs)
			if err != nil {
				return c.Status(500).SendString("Failed to delete")
			}
			err = userService.RemoveManyCommentIDs(story.CommentIDs)
			if err != nil {
				return c.Status(500).SendString("Failed to delete")
			}

			// --- remove related images in AWS S3 ---

			objects := make([]*s3.ObjectIdentifier, 0)

			for _, block := range story.Blocks {
				if block.Type == "image" {
					fileName := strings.Split(block.Data.File.URL, "amazonaws.com/")[1]
					objects = append(objects, &s3.ObjectIdentifier{Key: aws.String(fileName)})
				}
			}

			sess := myaws.ConnectAws()
			svc := s3.New(sess)
			bucketName := config.Config("BUCKET_NAME")
			_, err = svc.DeleteObjects(&s3.DeleteObjectsInput{Bucket: aws.String(bucketName), Delete: &s3.Delete{Objects: objects, Quiet: aws.Bool(true)}})
			if err != nil {
				fmt.Println(err)
				return c.SendStatus(500)
			}

			err = storyService.RemoveStory(storyOID)
			if err != nil {
				return c.Status(500).SendString("Failed to delete")
			}
		}

		// --- delete user's comments and related fields ---

		for _, commentOID := range *user.CommentIDs {
			comment, err := commentService.FindComment(commentOID)
			if err != nil {
				return c.Status(404).SendString("Comment not found")
			}
			if comment.CreatorID != userOID {
				return c.Status(400).SendString("You are not authorized.")
			}
			err = storyService.RemoveCommentID(comment.StoryID, commentOID)
			if err != nil {
				return c.Status(500).SendString("Failed to update")
			}
			err = commentService.RemoveComment(commentOID)
			if err != nil {
				return c.Status(500).SendString("Failed to delete")
			}
		}

		// --- delte from other users's followingIDs ---

		for _, followerOID := range *user.FollowerIDs {
			userService.RemoveFollowingID(followerOID, userOID)
		}

		// --- delete avatar photo ---

		bucketName := config.Config("BUCKET_NAME")

		sess := myaws.ConnectAws()
		svc := s3.New(sess)
		avatarFileName := strings.Split(user.AvatarURL, "amazonaws.com/")[1]
		_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucketName), Key: aws.String(avatarFileName)})
		if err != nil {
			fmt.Println(err)
			return c.Status(500).SendString("Failed to delete avatar image")
		}

		// --- delete user ---

		err = userService.RemoveAccount(userOID)
		if err != nil {
			return c.Status(500).SendString("Failed to delete account")
		}

		return c.SendStatus(204)
	}
}

func disbookmarkStory(userService user.Service) fiber.Handler {
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
		err = userService.DisbookmarkStory(userOID, storyOID)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}
		return c.SendStatus(200)
	}
}

func editPassword(userService user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type editPasswordInput struct {
			OriginalPass string `json:"originalPass"`
			FirstPass    string `json:"firstPass"`
			SecondPass   string `json:"secondPass"`
		}
		input := new(editPasswordInput)

		if err := c.BodyParser(input); err != nil {
			return c.SendStatus(400)
		}

		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}

		currentUser, err := userService.FindUserByID(userOID)
		if err != nil {
			return c.Status(404).SendString("User not found")
		}

		isValid := util.VerifyPassword(input.OriginalPass, currentUser.Password)
		if !isValid {
			return c.Status(400).SendString("You are not authorized")
		}
		if input.FirstPass != input.SecondPass {
			return c.Status(400).SendString("Verify your new password again.")
		}
		if len(input.FirstPass) < 6 {
			return c.Status(400).SendString("Password must be longer than 5 characters.")
		}
		newPass := util.HashPassword(input.FirstPass)
		err = userService.EditPassword(userOID, newPass)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}

		return c.SendStatus(200)
	}
}

func editAvatar(userService user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}

		file, err := c.FormFile("avatarUrl")
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(400)
		}
		oldURL := c.FormValue("oldAvatarUrl")
		oldFileName := strings.Split(oldURL, "amazonaws.com/")[1]

		f, err := file.Open()
		imageSrc, err := imaging.Decode(f)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}
		resizedImg := imaging.Resize(imageSrc, 200, 0, imaging.Lanczos)

		uuidWithHypen := uuid.New()
		uuid := strings.Replace(uuidWithHypen.String(), "-", "", -1)

		// ------
		// AWS S3
		// ------

		bucketName := config.Config("BUCKET_NAME")

		sess := myaws.ConnectAws()
		uploader := s3manager.NewUploader(sess)

		filename := uuid + file.Filename

		buf := new(bytes.Buffer)
		imaging.Encode(buf, resizedImg, imaging.JPEG)
		reader := bytes.NewReader(buf.Bytes())

		up, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucketName),
			ACL:    aws.String("public-read"),
			Key:    aws.String(filename),
			Body:   reader,
		})

		if err != nil {
			fmt.Println("Failed to upload file")
			return c.SendStatus(500)
		}

		if oldFileName != "blank-profile.webp" {
			svc := s3.New(sess)
			_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucketName), Key: aws.String(oldFileName)})
			if err != nil {
				fmt.Println(err)
				return c.SendStatus(500)
			}
		}

		err = userService.EditAvatar(userOID, up.Location)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}

		return c.Status(200).SendString(up.Location)
	}
}

func editBio(userService user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type editBioInput struct {
			Bio string `json:"bio"`
		}
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}

		input := new(editBioInput)
		if err := c.BodyParser(input); err != nil {
			return c.SendStatus(400)
		}
		err = userService.EditBio(userOID, input.Bio)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}
		return c.Status(200).SendString(input.Bio)
	}
}

func editUsername(userService user.Service) fiber.Handler {

	return func(c *fiber.Ctx) error {
		type editUsernameInput struct {
			Username string `json:"username"`
		}
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}

		input := new(editUsernameInput)
		if err := c.BodyParser(input); err != nil {
			return c.SendStatus(400)
		}

		err = userService.EditUsername(userOID, input.Username)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}

		return c.Status(200).SendString(input.Username)
	}
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
