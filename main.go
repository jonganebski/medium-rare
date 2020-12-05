package main

import (
	"home/jonganebski/github/fibersteps-server/config"
	"home/jonganebski/github/fibersteps-server/database"
	"home/jonganebski/github/fibersteps-server/middleware"
	"home/jonganebski/github/fibersteps-server/router"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/template/pug"
)

var port string = config.Config("PORT")

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}

	engine := pug.New("./views", ".pug")
	engine.AddFunc("publishBtn", func(path string) bool {
		if strings.Contains(path, "new-story") {
			return true
		}
		return false
	})
	engine.AddFunc("isMyStory", func(authorID, userID string) bool {
		if authorID == userID {
			return true
		}
		return false
	})
	engine.AddFunc("getStoryDate", func(createdAt int64) string {
		now := time.Now().Unix()
		lapse := now - createdAt
		oneDay := int64(24 * 60 * 60)
		if lapse < oneDay {
			return "today"
		}
		if lapse < 2*oneDay {
			return "yesterday"
		}
		if lapse < 3*oneDay {
			return "2 days ago"
		}
		return time.Unix(createdAt, 0).Format("January 2, 2006")
	})
	engine.AddFunc("grindBody", func(body string, targetLen int) string {
		if targetLen < len(body) {
			return body[:targetLen] + "..."
		}
		return body
	})

	app := fiber.New(fiber.Config{Views: engine})
	app.Static("/static", "./static")
	app.Static("/image", "./image")

	app.Use(helmet.New())
	app.Use(cors.New())
	app.Use(logger.New())

	app.Use(middleware.SpreadLocals)

	router.SetupRoutes(app)

	log.Fatal(app.Listen(":" + port))
}
