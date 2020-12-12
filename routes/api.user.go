package routes

import (
	"bytes"
	"fmt"
	myaws "home/jonganebski/github/medium-rare/aws"
	"home/jonganebski/github/medium-rare/config"
	"home/jonganebski/github/medium-rare/middleware"
	"home/jonganebski/github/medium-rare/package/comment"
	"home/jonganebski/github/medium-rare/package/story"
	"home/jonganebski/github/medium-rare/package/user"
	"home/jonganebski/github/medium-rare/util"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRouter has routes related with the user
func UserRouter(api fiber.Router, userService user.Service, storyService story.Service, commentService comment.Service) {
	// bookmark, disbookmark가 post와 delete여야할 이유가 없음.
	api.Post("/bookmark/:storyId", middleware.APIGuard, bookmarkStory(userService))
	api.Post("/follow/:authorId", middleware.APIGuard, follow(userService))
	api.Post("/unfollow/:authorId", middleware.APIGuard, unfollow(userService))
	api.Patch("/user/username", middleware.APIGuard, editUsername(userService))
	api.Patch("/user/bio", middleware.APIGuard, editBio(userService))
	api.Patch("/user/avatar", middleware.APIGuard, editAvatar(userService))
	api.Patch("/user/password", middleware.APIGuard, editPassword(userService))
	api.Delete("/bookmark/:storyId", middleware.APIGuard, disbookmarkStory(userService))
	api.Delete("/user", middleware.APIGuard, removeAccount(userService, storyService, commentService))
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
			return c.Status(403).SendString("Wrong password")
		}

		// --- delete user's stories and related fields ---
		for _, storyOID := range *user.StoryIDs {
			// find story
			story, err := storyService.FindStoryByID(storyOID)
			if err != nil {
				return c.Status(404).SendString("Story not found")
			}
			// check current user is the author of the story
			if story.CreatorID != userOID {
				return c.Status(400).SendString("You are not authorized.")
			}
			// remove storyID from other users' liked story IDs field
			err = userService.RemoveManyLikedStoryIDs(storyOID)
			if err != nil {
				return c.Status(500).SendString("Failed to update")
			}
			// remove storyID from other users' bookmarked story IDs field
			err = userService.RemoveManySavedStoryIDs(storyOID)
			if err != nil {
				return c.Status(500).SendString("Failed to update")
			}
			// remove the story's comments (comment document itself)
			err = commentService.RemoveComments(story.CommentIDs)
			if err != nil {
				return c.Status(500).SendString("Failed to delete")
			}
			// remove commentID from each user's commentIDs field
			for _, commentID := range *story.CommentIDs {
				err = userService.RemoveManyUsersCommentID(commentID)
				if err != nil {
					return c.Status(500).SendString("Failed to delete")
				}
			}

			// remove related images in AWS S3
			objects := make([]*s3.ObjectIdentifier, 0)

			for _, block := range story.Blocks {
				// make a slice of image file names of the story
				if block.Type == "image" {
					fileName := strings.Split(block.Data.File.URL, "amazonaws.com/")[1]
					objects = append(objects, &s3.ObjectIdentifier{Key: aws.String(fileName)})
				}
			}

			if len(objects) != 0 {
				// in case of deleting many objects, aws throws error when the slice of objects is empty
				sess := myaws.ConnectAws()
				svc := s3.New(sess)
				bucketName := config.Config("BUCKET_NAME")
				_, err = svc.DeleteObjects(&s3.DeleteObjectsInput{Bucket: aws.String(bucketName), Delete: &s3.Delete{Objects: objects, Quiet: aws.Bool(true)}})
				if err != nil {
					fmt.Println(err)
					return c.SendStatus(500)
				}
			}
			// remove story document itself
			err = storyService.RemoveStory(storyOID)
			if err != nil {
				return c.Status(500).SendString("Failed to delete")
			}
		}

		// --- find user again as commentIDs can be changed when story was being removed (if user commented his/her own story) ---
		user, err = userService.FindUserByID(userOID)
		if err != nil {
			return c.Status(404).SendString("User not found")
		}

		// --- delete user's comments and related fields ---
		for _, commentOID := range *user.CommentIDs {
			// find comment
			comment, err := commentService.FindComment(commentOID)
			if err != nil {
				return c.Status(404).SendString("Comment not found")
			}
			// check the comment belongs to current user
			if comment.CreatorID != userOID {
				return c.Status(400).SendString("You are not authorized.")
			}
			// remove commentID from any other stories' commentIDs field
			err = storyService.RemoveCommentID(comment.StoryID, commentOID)
			if err != nil {
				return c.Status(500).SendString("Failed to update")
			}
			// remove comment document itself
			err = commentService.RemoveComment(commentOID)
			if err != nil {
				return c.Status(500).SendString("Failed to delete")
			}
		}

		// --- remove userID from other users's followingIDs ---
		for _, followerOID := range *user.FollowerIDs {
			userService.RemoveFollowingID(followerOID, userOID)
		}

		// --- remove userID from other users's followerIDs ---
		for _, followingOID := range *user.FollowingIDs {
			userService.RemoveFollowerID(followingOID, userOID)
		}

		// --- remove avatar photo in aws-s3 ---
		bucketName := config.Config("BUCKET_NAME")

		sess := myaws.ConnectAws()
		avatarFileName := strings.Split(user.AvatarURL, "amazonaws.com/")[1]
		if avatarFileName != "blank-profile.webp" {
			svc := s3.New(sess)
			_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucketName), Key: aws.String(avatarFileName)})
			if err != nil {
				fmt.Println(err)
				return c.SendStatus(500)
			}
		}

		// --- delete user document itself ---
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
		// remove storyID from current user's savedStoryIDs field
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

		// --- find current user ---
		currentUser, err := userService.FindUserByID(userOID)
		if err != nil {
			return c.Status(404).SendString("User not found")
		}

		// --- verify password ---
		isValid := util.VerifyPassword(input.OriginalPass, currentUser.Password)
		if !isValid {
			return c.Status(400).SendString("You are not authorized")
		}

		// --- validate new password ---
		if input.FirstPass != input.SecondPass {
			return c.Status(400).SendString("Verify your new password again.")
		}
		if len(input.FirstPass) < 6 {
			return c.Status(400).SendString("Password must be longer than 5 characters.")
		}

		// --- hash new password and update user's password field ---
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

		// --- open and decode image file ---
		f, err := file.Open()
		imageSrc, err := imaging.Decode(f)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}

		// --- resize image ---
		resizedImg := imaging.Resize(imageSrc, 200, 0, imaging.Lanczos)

		// --- name image file with uuid ---
		uuidWithHypen := uuid.New()
		uuid := strings.Replace(uuidWithHypen.String(), "-", "", -1)
		filename := uuid + file.Filename

		// ------
		// AWS S3
		// ------

		bucketName := config.Config("BUCKET_NAME")
		sess := myaws.ConnectAws()
		uploader := s3manager.NewUploader(sess)

		// --- encode image and make reader ---
		buf := new(bytes.Buffer)
		imaging.Encode(buf, resizedImg, imaging.JPEG)
		reader := bytes.NewReader(buf.Bytes())

		// --- upload avatar image to aws-s3 ---
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

		// --- remove previous avatar image ---
		if oldFileName != "blank-profile.webp" {
			svc := s3.New(sess)
			_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucketName), Key: aws.String(oldFileName)})
			if err != nil {
				fmt.Println(err)
				return c.SendStatus(500)
			}
		}

		// --- update user's avatar url ---
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
		input := new(editBioInput)
		if err := c.BodyParser(input); err != nil {
			return c.SendStatus(400)
		}
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}

		// --- update user's bio field ---
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
		input := new(editUsernameInput)
		if err := c.BodyParser(input); err != nil {
			return c.SendStatus(400)
		}
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}

		// --- update user's username field ---
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

		// --- remove current user's ID from target user's FollowingIDs field ---
		err = userService.RemoveFollowerID(authorOID, userOID)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}

		// --- remove target user's ID from current user's FollowerIDs field ---
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

		// --- find target user ---
		author, err := userService.FindUserByID(authorOID)
		if err != nil {
			return c.Status(404).SendString("Author not found")
		}

		// --- check current user is already following target user ---
		for _, followerID := range *author.FollowerIDs {
			if followerID == userOID {
				return c.Status(400).SendString("You are already following this user.")
			}
		}

		// --- add target user's ID into current user's FollowingIDs field ---
		err = userService.AddFollowingID(userOID, authorOID)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}

		// --- add current user's ID into target user's FollowerIDs field ---
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

		// --- add story's ID into currnet user's SavedStoryIDs field ---
		err = userService.BookmarkStory(userOID, storyOID)
		if err != nil {
			return err
		}

		return c.SendStatus(200)
	}
}
