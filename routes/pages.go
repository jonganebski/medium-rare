package routes

import (
	"fmt"
	"home/jonganebski/github/medium-rare/helper"
	"home/jonganebski/github/medium-rare/middleware"
	"home/jonganebski/github/medium-rare/model"
	"home/jonganebski/github/medium-rare/package/story"
	"home/jonganebski/github/medium-rare/package/user"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type storyCardOutput struct {
	StoryID        string `json:"storyId"`
	AuthorID       string `json:"authorId"`
	AuthorUsername string `json:"authorUsername"`
	CreatedAt      int64  `json:"createdAt"`
	UpdatedAt      int64  `json:"updatedAt"`
	Header         string `json:"header"`
	Body           string `json:"body"`
	CoverImgURL    string `json:"coverImgUrl"`
	ReadTime       string `json:"readTime"`
	Ranking        int    `json:"ranking,omitempty"`
	IsPublished    bool   `json:"isPublished"`
}

type followerOutput struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	AvatarURL    string `json:"avatarUrl"`
	AmIFollowing bool   `json:"amIFollowing"`
	IsMe         bool   `json:"isMe"`
}

// PageRouter has the routes where renders webpage
func PageRouter(root fiber.Router, userService user.Service, storyService story.Service) {
	root.Get("/", homepage(userService, storyService))
	root.Get("/read-story/:storyId", readStory(userService, storyService))
	root.Get("/followers/:userId", seeFollowers(userService))
	root.Get("/user-home/:userId", userHome(userService, storyService))
	root.Get("/new-story", middleware.Protected, newStory(userService))
	root.Get("/edit-story/:storyId", middleware.Protected, editStoryPage(userService, storyService))

	me := root.Group("/me", middleware.Protected)
	me.Get("/bookmarks", myBookmarks(userService, storyService))
	me.Get("/following", seeFollowings(userService))
	me.Get("/settings", settings(userService))
	me.Get("/stories", myStories(userService, storyService))
}

func myStories(userService user.Service, storyService story.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}
		// --- find current user ---
		user, err := userService.FindUserByID(userOID)
		if err != nil {
			return c.Status(404).SendString("User not found")
		}
		// --- find current user's stories ---
		stories, err := storyService.FindStories(user.StoryIDs)
		if err != nil {
			return c.Status(404).SendString("Stories not found")
		}
		// --- make output for pug partial storyCard ---
		storyCards, err := composeStoryCardOutput(*stories, userService)
		if err != nil {
			// 여기서 저자는 곧 유저이므로 찾을 필요가 없긴 하다. 나중에 고쳐볼 것.
			return c.Status(404).SendString("Author not found")
		}
		return c.Render("myStories", fiber.Map{
			"path":        c.Path(),
			"userId":      c.Locals("userId"),
			"currentUser": user,
			"storyCards":  storyCards,
		}, "layout/main")
	}
}

func settings(userService user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}
		// --- find current user ---
		currentUser, err := userService.FindUserByID(userOID)
		if err != nil {
			return c.Status(404).SendString("User not found")
		}
		return c.Render("settings", fiber.Map{
			"path":           c.Path(),
			"userId":         c.Locals("userId"),
			"currentUser":    currentUser,
			"followerCount":  len(*currentUser.FollowerIDs),
			"followingCount": len(*currentUser.FollowingIDs),
		}, "layout/main")
	}
}

func seeFollowings(userService user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}
		// --- find current user ---
		currentUser, err := userService.FindUserByID(userOID)
		if err != nil {
			return c.Status(404).SendString("User not found")
		}
		// --- find people current user is following ---
		followings, err := userService.FindUsers(currentUser.FollowingIDs)
		if err != nil {
			return c.Status(404).SendString("Followings not found")
		}
		return c.Render("following", fiber.Map{
			"path":        c.Path(),
			"userId":      c.Locals("userId"),
			"currentUser": currentUser,
			"followings":  followings,
		}, "layout/main")
	}
}

func myBookmarks(userService user.Service, storyService story.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}
		// --- find current user ---
		currentUser, err := userService.FindUserByID(userOID)
		if err != nil {
			return c.Status(404).SendString("User not found")
		}
		// --- find current user's bookmarked stories ---
		stories, err := storyService.FindStories(currentUser.SavedStoryIDs)
		if err != nil {
			return c.Status(404).SendString("Stories not found")
		}
		// --- make output for pug partial storyCard ---
		storyCards, err := composeStoryCardOutput(*stories, userService)
		if err != nil {
			return c.Status(404).SendString("Author not found")
		}
		return c.Render("bookmarks", fiber.Map{
			"path":        c.Path(),
			"userId":      c.Locals("userId"),
			"currentUser": currentUser,
			"storyCards":  storyCards,
		}, "layout/main")
	}
}

func userHome(userService user.Service, storyService story.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		targetUserID := c.Params("userId")
		targetUserOID, err := primitive.ObjectIDFromHex(targetUserID)
		if err != nil {
			return c.SendStatus(500)
		}
		// --- find user of the page ---
		targetUser, err := userService.FindUserByID(targetUserOID)
		if err != nil {
			return c.Status(404).SendString("User not found")
		}
		// --- find page user's stories ---
		stories, err := storyService.FindStories(targetUser.StoryIDs)
		if err != nil {
			return c.Status(404).SendString("Stories not found")
		}
		// --- make output for pug partial storyCard ---
		storyCards, err := composeStoryCardOutput(*stories, userService)
		if err != nil {
			return c.Status(404).SendString("Author not found")
		}

		var publishedStoryCards []storyCardOutput
		for _, storyCard := range *storyCards {
			if storyCard.IsPublished {
				publishedStoryCards = append(publishedStoryCards, storyCard)
			}
		}

		currentUser := new(model.User)
		isFollowingTargetUser := false
		if c.Locals("userId") != nil {
			// if current user is signed-in
			userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
			if err != nil {
				return c.SendStatus(500)
			}
			// find current user
			currentUser, err = userService.FindUserByID(userOID)
			// check relationship between current user and user of the page
			for _, followingID := range *currentUser.FollowingIDs {
				if followingID == targetUserOID {
					isFollowingTargetUser = true
					break
				}
			}
		}
		return c.Render("user-home", fiber.Map{
			"path":                  c.Path(),
			"userId":                c.Locals("userId"),
			"currentUser":           currentUser,
			"targetUser":            targetUser,
			"storyCards":            publishedStoryCards,
			"storyCardsCount":       len(publishedStoryCards),
			"isFollowingTargetUser": isFollowingTargetUser,
		}, "layout/main")
	}
}

func seeFollowers(userService user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		targetUserID := c.Params("userId")
		targetUserOID, err := primitive.ObjectIDFromHex(targetUserID)
		if err != nil {
			return c.SendStatus(500)
		}
		// --- find user of the page ---
		targetUser, err := userService.FindUserByID(targetUserOID)
		if err != nil {
			return c.Status(404).SendString("This user is not found")
		}

		currentUser := new(model.User)
		outputItem := new(followerOutput)
		output := make([]followerOutput, 0)
		isFollowingTargetUser := false

		if c.Locals("userId") != nil {
			// --- if current user is signed-in ---
			currentUserOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
			if err != nil {
				return c.SendStatus(500)
			}
			// find current user
			currentUser, err = userService.FindUserByID(currentUserOID)
			if err != nil {
				return c.Status(404).SendString("User not found")
			}
			// make output struct for pug partial userCard
			for _, followerID := range *targetUser.FollowerIDs {
				follower, err := userService.FindUserByID(followerID)
				if err != nil {
					return c.SendStatus(404)
				}
				outputItem.ID = followerID.Hex()
				outputItem.Username = follower.Username
				outputItem.AvatarURL = follower.AvatarURL
				outputItem.IsMe = (currentUserOID == followerID)
				outputItem.AmIFollowing = false
				// check relationship between current user and user of the page
				if followerID == currentUserOID {
					isFollowingTargetUser = true
				}
				// check relationship between current user and followers of user of the page
				for _, followingID := range *currentUser.FollowingIDs {
					if followingID == followerID {
						outputItem.AmIFollowing = true
						break
					}
				}
				output = append(output, *outputItem)
			}
		} else {
			// --- if current user is not signed-in ---
			// make output struct for pug partial userCard
			for _, followerID := range *targetUser.FollowerIDs {
				follower, err := userService.FindUserByID(followerID)
				if err != nil {
					return c.SendStatus(404)
				}
				outputItem.ID = followerID.Hex()
				outputItem.Username = follower.Username
				outputItem.AvatarURL = follower.AvatarURL
				outputItem.IsMe = false
				outputItem.AmIFollowing = false
				output = append(output, *outputItem)
			}
		}
		return c.Render("followers", fiber.Map{
			"path":                  c.Path(),
			"userId":                c.Locals("userId"),
			"currentUser":           currentUser,
			"targetUser":            targetUser,
			"followersOutput":       output,
			"isFollowingTargetUser": isFollowingTargetUser,
		}, "layout/main")
	}
}

func editStoryPage(userService user.Service, storyService story.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		storyID := c.Params("storyId")
		storyOID, err := primitive.ObjectIDFromHex(storyID)
		if err != nil {
			return c.SendStatus(500)
		}
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			fmt.Println("error at conversion")
			return c.SendStatus(500)
		}
		// --- find story ---
		story, err := storyService.FindStoryByID(storyOID)
		if err != nil {
			return c.Status(404).SendString("Story not found")
		}
		// --- check the story belongs to current user ---
		if userOID != story.CreatorID {
			fmt.Println("You are not authorized")
			c.Redirect("/")
		}
		// --- find current user ---
		currentUser, err := userService.FindUserByID(userOID)
		if err != nil {
			return c.Status(404).SendString("User not found")
		}
		return c.Render("editStory", fiber.Map{
			"path":        c.Path(),
			"userId":      c.Locals("userId"),
			"currentUser": currentUser,
			"story":       story,
		}, "layout/main")
	}
}

func readStory(userService user.Service, storyService story.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		storyID := c.Params("storyId")
		storyOID, err := primitive.ObjectIDFromHex(storyID)
		if err != nil {
			return c.SendStatus(500)
		}
		// --- increase view count of the story ---
		story, err := storyService.IncreaseViewCount(storyOID)
		if err != nil {
			return c.Status(404).SendString("Story not found")
		}
		// --- find author user of the story ---
		author, err := userService.FindUserByID(story.CreatorID)
		if err != nil {
			return c.Status(404).SendString("Author not found")
		}

		didLiked := false
		bookmarked := false
		isFollowing := false
		currentUser := new(model.User)
		if c.Locals("userId") != nil {
			// --- if currnet user is signed in ---
			userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
			if err != nil {
				return c.SendStatus(500)
			}
			// find current user
			currentUser, err = userService.FindUserByID(userOID)
			if err != nil {
				c.ClearCookie()
				return c.Redirect("/")
			}
			// did currnet user liked this story?
			for _, likedStoryID := range *currentUser.LikedStoryIDs {
				if likedStoryID == storyOID {
					didLiked = true
					break
				}
			}
			// did current user bookmarked this story?
			for _, savedStoryID := range *currentUser.SavedStoryIDs {
				if savedStoryID == storyOID {
					bookmarked = true
					break
				}
			}
			// is current user following author of the story
			for _, followerID := range *author.FollowerIDs {
				if followerID == userOID {
					isFollowing = true
					break
				}
			}
		}
		return c.Render("readStory", fiber.Map{
			"path":        c.Path(),
			"userId":      c.Locals("userId"),
			"currentUser": currentUser,
			"story":       story,
			"author":      author,
			"didLiked":    didLiked,
			"bookmarked":  bookmarked,
			"isFollowing": isFollowing,
		}, "layout/main")
	}
}

func newStory(userService user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}
		// --- find current user ---
		user, err := userService.FindUserByID(userOID)
		if err != nil {
			return c.Redirect("/")
		}
		return c.Render("newStory", fiber.Map{
			"path":        c.Path(),
			"userId":      c.Locals("userId"),
			"currentUser": user,
		}, "layout/main")
	}
}

func homepage(userService user.Service, storyService story.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentUser := new(model.User)
		if c.Locals("userId") != nil {
			// --- if current user is signed-in ---
			userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
			if err != nil {
				fmt.Println(err)
				return c.SendStatus(500)
			}
			// find current user
			currentUser, err = userService.FindUserByID(userOID)
			if err != nil {
				c.ClearCookie()
				return c.Redirect("/")
			}
		}
		// --- find stories by time ---
		stories, err := storyService.FindRecentStories()
		if err != nil {
			return c.SendStatus(404)
		}
		// --- make output for pug partial storyCard ---
		recentStories, err := composeStoryCardOutput(*stories, userService)
		if err != nil {
			return c.SendStatus(404)
		}
		// --- find stories picked by editor ---
		stories, err = storyService.FindPickedStories()
		if err != nil {
			return c.SendStatus(404)
		}
		// --- make output for pug partial storyCard ---
		pickedStories, err := composeStoryCardOutput(*stories, userService)
		if err != nil {
			return c.SendStatus(404)
		}
		// --- find stories by view count ---
		stories, err = storyService.FindPopularStories()
		if err != nil {
			return c.SendStatus(404)
		}
		// --- make output for pug partial storyCard ---
		popularStories, err := composeStoryCardOutput(*stories, userService)
		if err != nil {
			return c.SendStatus(404)
		}
		// --- split stories picked by editor ---
		// this can be unnecessary if I fix structure of grid in pug file
		editorsPickR := new(storyCardOutput)
		editorsPickC := make([]storyCardOutput, 0)
		editorsPickL := new(storyCardOutput)

		for i, output := range *pickedStories {
			if i == 0 {
				editorsPickR = &output
			}
			if 0 < i && i < 4 {
				editorsPickC = append(editorsPickC, output)
			}
			if i == 4 {
				editorsPickL = &output
			}
		}

		return c.Render("home", fiber.Map{
			"path":           c.Path(),
			"userId":         c.Locals("userId"),
			"currentUser":    currentUser,
			"editorsPickR":   editorsPickR,
			"editorsPickC":   editorsPickC,
			"editorsPickL":   editorsPickL,
			"recentStories":  recentStories,
			"popularStories": popularStories,
		}, "layout/main")
	}
}

func composeStoryCardOutput(stories []model.Story, userService user.Service) (*[]storyCardOutput, error) {
	storyCard := new(storyCardOutput)
	storyCards := make([]storyCardOutput, 0)
	for _, story := range stories {
		// --- find body & coverImgUrl & compute readTime ---
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
		// --- find author ---
		author, err := userService.FindUserByID(story.CreatorID)
		if err != nil {
			return nil, err
		}
		// --- build outputItem and append to output ---
		storyCard.AuthorUsername = author.Username
		storyCard.StoryID = story.ID
		storyCard.AuthorID = author.ID
		storyCard.Header = story.Blocks[0].Data.Text
		storyCard.Body = body
		storyCard.CreatedAt = story.CreatedAt
		storyCard.UpdatedAt = story.UpdatedAt
		storyCard.CoverImgURL = coverImgURL
		storyCard.ReadTime = readTimeText
		storyCard.IsPublished = story.IsPublished
		storyCards = append(storyCards, *storyCard)
	}
	return &storyCards, nil
}
