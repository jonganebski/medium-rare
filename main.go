package main

import (
	"home/jonganebski/github/medium-rare/config"
	"home/jonganebski/github/medium-rare/database"
	"home/jonganebski/github/medium-rare/helper"
	"home/jonganebski/github/medium-rare/middleware"
	"home/jonganebski/github/medium-rare/router"
	"log"

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
	engine.AddFunc("publishBtn", helper.IsPublishButton)
	engine.AddFunc("isMyStory", helper.IsMyStory)
	engine.AddFunc("getStoryDate", helper.GetStoryPostDate)
	engine.AddFunc("grindBody", helper.GrindBody)
	engine.AddFunc("getSliceLen", helper.GetSliceLen)
	engine.AddFunc("SortByUpdatedAt", helper.SortByUpdatedAt)
	engine.AddFunc("getYear", helper.GetYear)

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
