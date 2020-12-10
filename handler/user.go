package handler

import (
	"bytes"
	"fmt"
	myaws "home/jonganebski/github/medium-rare/aws"
	"home/jonganebski/github/medium-rare/config"
	"home/jonganebski/github/medium-rare/database"
	"home/jonganebski/github/medium-rare/helper"
	"home/jonganebski/github/medium-rare/model"
	"home/jonganebski/github/medium-rare/util"
	"image"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	// --- find author ---

	author := new(model.User)
	filter := bson.D{{Key: "_id", Value: authorOID}}
	singleResult := userCollection.FindOne(c.Context(), filter)
	if singleResult.Err() != nil {
		return c.SendStatus(404)
	}
	singleResult.Decode(author)

	// --- check the user is following the author already or not

	for _, followerID := range *author.FollowerIDs {
		if followerID == userOID {
			return c.SendStatus(400)
		}
	}

	// --- add author's userID into current user's FollowingIDs ---

	filter = bson.D{{Key: "_id", Value: userOID}}
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

	// --- remove author's userID from current user's FollowingIDs ---

	filter := bson.D{{Key: "_id", Value: userOID}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "followingIds", Value: authorOID}}}}
	updateResult := userCollection.FindOneAndUpdate(c.Context(), filter, update)
	if updateResult.Err() != nil {
		return c.SendStatus(404)
	}

	// --- remove current user's userID from author's FollowingIDs ---

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

	return c.Render("settings", fiber.Map{"path": c.Path(), "userId": c.Locals("userId"), "userAvatarUrl": user.AvatarURL, "username": user.Username, "userEmail": user.Email, "bio": user.Bio, "followerCount": len(*user.FollowerIDs), "followingCount": len(*user.FollowingIDs)}, "layout/main")
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

	userCollection := mg.Db.Collection(UserCollection)

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
	imageSrc, _, err := image.Decode(f)
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

	svc := s3.New(sess)
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucketName), Key: aws.String(oldFileName)})
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(500)
	}

	user := new(model.User)
	filter := bson.D{{Key: "_id", Value: userOID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "avatarUrl", Value: up.Location}}}}
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

// EditPassword changes current user's password
func EditPassword(c *fiber.Ctx) error {

	type editPasswordInput struct {
		OriginalPass string `json:"originalPass"`
		FirstPass    string `json:"firstPass"`
		SecondPass   string `json:"secondPass"`
	}
	userCollection := mg.Db.Collection(UserCollection)

	input := new(editPasswordInput)

	if err := c.BodyParser(input); err != nil {
		return c.SendStatus(400)
	}

	userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
	if err != nil {
		return c.SendStatus(500)
	}

	user := new(model.User)
	filter := bson.D{{Key: "_id", Value: userOID}}
	singleResult := userCollection.FindOne(c.Context(), filter)
	if singleResult.Err() != nil {
		return c.SendStatus(404)
	}
	singleResult.Decode(user)

	isValid := util.VerifyPassword(input.OriginalPass, user.Password)
	if !isValid {
		return c.SendStatus(400)
	}
	if input.FirstPass != input.SecondPass {
		return c.SendStatus(400)
	}
	if len(input.FirstPass) < 6 {
		return c.SendStatus(400)
	}

	newPass := util.HashPassword(input.FirstPass)

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: newPass}}}}
	updateResult, err := userCollection.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return c.SendStatus(500)
	}
	if updateResult.ModifiedCount == 0 {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)
}

// SeeFollowers renders a page where shows target user's followers
func SeeFollowers(c *fiber.Ctx) error {

	type followerOutput struct {
		ID           string `json:"id"`
		Username     string `json:"username"`
		AvatarURL    string `json:"avatarUrl"`
		AmIFollowing bool   `json:"amIFollowing"`
		IsMe         bool   `json:"isMe"`
	}

	userCollection := mg.Db.Collection(UserCollection)

	targetUserID := c.Params("userId")
	targetUserOID, err := primitive.ObjectIDFromHex(targetUserID)
	if err != nil {
		return c.SendStatus(500)
	}

	targetUser := new(model.User)
	filter := bson.D{{Key: "_id", Value: targetUserOID}}
	singleResult := userCollection.FindOne(c.Context(), filter)
	if singleResult.Err() != nil {
		return c.SendStatus(404)
	}
	singleResult.Decode(targetUser)

	user := new(model.User)
	follower := new(followerOutput)
	followers := make([]followerOutput, 0)

	if c.Locals("userId") != nil {
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}

		filter = bson.D{{Key: "_id", Value: userOID}}
		singleResult = userCollection.FindOne(c.Context(), filter)
		if singleResult.Err() != nil {
			return c.SendStatus(404)
		}
		singleResult.Decode(user)
		for _, followerID := range *targetUser.FollowerIDs {
			filter := bson.D{{Key: "_id", Value: followerID}}
			singleResult := userCollection.FindOne(c.Context(), filter)
			if singleResult.Err() != nil {
				return c.SendStatus(404)
			}
			singleResult.Decode(follower)
			follower.ID = followerID.Hex()
			follower.IsMe = (userOID == followerID)
			follower.AmIFollowing = false
			for _, followingID := range *user.FollowingIDs {
				if followingID == followerID {
					follower.AmIFollowing = true
					break
				}
			}
			followers = append(followers, *follower)
		}
	} else {
		for _, followerID := range *targetUser.FollowerIDs {
			filter := bson.D{{Key: "_id", Value: followerID}}
			singleResult := userCollection.FindOne(c.Context(), filter)
			if singleResult.Err() != nil {
				return c.SendStatus(404)
			}
			singleResult.Decode(follower)
			follower.ID = followerID.Hex()
			follower.IsMe = false
			follower.AmIFollowing = false
			followers = append(followers, *follower)
		}
	}

	return c.Render("followers",
		fiber.Map{
			"path":          c.Path(),
			"userId":        c.Locals("userId"),
			"userAvatarUrl": user.AvatarURL,
			"username":      user.Username,
			"userEmail":     user.Email,
			"targetUser":    targetUser,
			"followers":     followers,
		},
		"layout/main")
}

// SeeFollowings renders a page where shows users target user is following
func SeeFollowings(c *fiber.Ctx) error {

	type followingOutput struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		AvatarURL string `json:"avatarUrl"`
	}

	userCollection := mg.Db.Collection(UserCollection)

	user := new(model.User)
	if c.Locals("userId") != nil {
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}

		filter := bson.D{{Key: "_id", Value: userOID}}
		singleResult := userCollection.FindOne(c.Context(), filter)
		if singleResult.Err() != nil {
			return c.SendStatus(404)
		}
		singleResult.Decode(user)
	}

	following := new(followingOutput)
	followings := make([]*followingOutput, 0)

	for _, followingID := range *user.FollowingIDs {
		filter := bson.D{{Key: "_id", Value: followingID}}
		singleResult := userCollection.FindOne(c.Context(), filter)
		if singleResult.Err() != nil {
			return c.SendStatus(404)
		}
		singleResult.Decode(following)
		following.ID = followingID.Hex()
		followings = append(followings, following)
	}

	return c.Render("following", fiber.Map{"path": c.Path(), "userId": c.Locals("userId"), "username": user.Username, "userAvatarUrl": user.AvatarURL, "userEmail": user.Email, "followingIDs": user.FollowingIDs, "followings": followings}, "layout/main")
}

// UserHome renders target user's page
func UserHome(c *fiber.Ctx) error {

	type storyCardOutput struct {
		StoryID        string `json:"storyId"`
		AuthorID       string `json:"authorId"`
		AuthorUsername string `json:"authorUsername"`
		CreatedAt      int64  `json:"createdAt"`
		Header         string `json:"header"`
		Body           string `json:"body"`
		CoverImgURL    string `json:"coverImgUrl"`
		ReadTime       string `json:"readTime"`
		Ranking        int    `json:"ranking,omitempty"`
	}

	userCollection := mg.Db.Collection(UserCollection)
	storyCollection := mg.Db.Collection(StoryCollection)

	outputItem := new(storyCardOutput)
	output := make([]storyCardOutput, 0)

	targetUserID := c.Params("userId")
	targetUserOID, err := primitive.ObjectIDFromHex(targetUserID)
	if err != nil {
		return c.SendStatus(500)
	}
	targetUser := new(model.User)
	filter := bson.D{{Key: "_id", Value: targetUserOID}}
	singleResult := userCollection.FindOne(c.Context(), filter)
	if singleResult.Err() != nil {
		return c.SendStatus(404)
	}
	singleResult.Decode(targetUser)

	stories := make([]model.Story, 0)
	filter = bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: targetUser.StoryIDs}}}}
	cursor, err := storyCollection.Find(c.Context(), filter)
	if err != nil {
		return c.SendStatus(500)
	}
	if err = cursor.All(c.Context(), &stories); err != nil {
		return c.SendStatus(500)
	}

	for _, story := range stories {

		body := ""
		coverImgURL := ""
		totalText := ""
		for _, block := range story.Blocks {
			if block.Type == "paragraph" {
				totalText += block.Data.Text
				if body == "" {
					body = block.Data.Text
				}
			}
			if block.Type == "image" && coverImgURL == "" {
				coverImgURL = block.Data.File.URL
			}
			if block.Type == "code" {
				totalText += block.Data.Code
			}
		}
		readTimeText := helper.ComputeReadTime(totalText)

		outputItem.AuthorUsername = targetUser.Username
		outputItem.StoryID = story.ID
		outputItem.Header = story.Blocks[0].Data.Text
		outputItem.Body = body
		outputItem.CreatedAt = story.CreatedAt
		outputItem.CoverImgURL = coverImgURL
		outputItem.ReadTime = readTimeText
		output = append(output, *outputItem)
	}

	user := new(model.User)
	if c.Locals("userId") != nil {
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}

		filter = bson.D{{Key: "_id", Value: userOID}}
		singleResult = userCollection.FindOne(c.Context(), filter)
		if singleResult.Err() != nil {
			return c.SendStatus(404)
		}
		singleResult.Decode(user)
	}

	return c.Render("user-home",
		fiber.Map{
			"path":          c.Path(),
			"userId":        c.Locals("userId"),
			"username":      user.Username,
			"userAvatarUrl": user.AvatarURL,
			"userEmail":     user.Email,
			"targetUser":    targetUser,
			"output":        output,
		},
		"layout/main")
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
