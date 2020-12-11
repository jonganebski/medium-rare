package routes

import (
	"fmt"
	"home/jonganebski/github/medium-rare/helper"
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
	Header         string `json:"header"`
	Body           string `json:"body"`
	CoverImgURL    string `json:"coverImgUrl"`
	ReadTime       string `json:"readTime"`
	Ranking        int    `json:"ranking,omitempty"`
}

// PageRouter has the routes where renders webpage
func PageRouter(app fiber.Router, userService user.Service, storyService story.Service) {
	app.Get("/", homepage(userService, storyService))
}

func homepage(userService user.Service, storyService story.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currnetUser := new(model.User)
		if c.Locals("userId") != nil {
			userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
			if err != nil {
				fmt.Println(err)
				return c.SendStatus(500)
			}
			currnetUser, err = userService.FindUserByID(userOID)
			if err != nil {
				c.ClearCookie()
				return c.Redirect("/")
			}
		}
		stories, err := storyService.FindRecentStories()
		if err != nil {
			return c.SendStatus(404)
		}
		recentStories, err := composeStoryCardOutput(*stories, userService)
		if err != nil {
			return c.SendStatus(404)
		}

		stories, err = storyService.FindPickedStories()
		if err != nil {
			return c.SendStatus(404)
		}
		pickedStories, err := composeStoryCardOutput(*stories, userService)
		if err != nil {
			return c.SendStatus(404)
		}

		stories, err = storyService.FindPopularStories()
		if err != nil {
			return c.SendStatus(404)
		}
		popularStories, err := composeStoryCardOutput(*stories, userService)
		if err != nil {
			return c.SendStatus(404)
		}

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
			"currentUser":    currnetUser,
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

		// find body & coverImgUrl & compute readTime

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

		// find author

		author, err := userService.FindUserByID(story.CreatorID)
		if err != nil {
			return nil, err
		}

		// build outputItem and append to output

		storyCard.AuthorUsername = author.Username
		storyCard.StoryID = story.ID
		storyCard.AuthorID = author.ID
		storyCard.Header = story.Blocks[0].Data.Text
		storyCard.Body = body
		storyCard.CreatedAt = story.CreatedAt
		storyCard.CoverImgURL = coverImgURL
		storyCard.ReadTime = readTimeText
		storyCards = append(storyCards, *storyCard)
	}
	return &storyCards, nil
}
