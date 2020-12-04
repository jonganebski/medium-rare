package main

import (
	"home/jonganebski/github/fibersteps-server/config"
	"home/jonganebski/github/fibersteps-server/database"
	"home/jonganebski/github/fibersteps-server/router"
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

	app := fiber.New(fiber.Config{Views: engine})
	app.Static("/static", "./static")

	app.Use(helmet.New())
	app.Use(cors.New())
	app.Use(logger.New())

	router.SetupRoutes(app)

	log.Fatal(app.Listen(":" + port))
}
