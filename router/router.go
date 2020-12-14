package router

import (
	myaws "home/jonganebski/github/medium-rare/aws"
	"home/jonganebski/github/medium-rare/config"
	"home/jonganebski/github/medium-rare/database"
	"home/jonganebski/github/medium-rare/package/comment"
	"home/jonganebski/github/medium-rare/package/photo"
	"home/jonganebski/github/medium-rare/package/story"
	"home/jonganebski/github/medium-rare/package/user"
	"home/jonganebski/github/medium-rare/routes"

	"github.com/gofiber/fiber/v2"
)

var mg = &database.Mongo

// SetupRoutes sets up the routes
func SetupRoutes(app *fiber.App) {
	// --- mongodb collections & aws s3 ---
	userCollection := mg.Db.Collection(config.Config("COLLECTION_USER"))
	storyCollection := mg.Db.Collection(config.Config("COLLECTION_STORY"))
	commentCollection := mg.Db.Collection(config.Config("COLLECTION_COMMENT"))
	sess := myaws.ConnectAws()

	// --- initialize repositories ---
	userRepo := user.NewRepo(userCollection)
	storyRepo := story.NewRepo(storyCollection)
	commentRepo := comment.NewRepo(commentCollection)
	uploadRepo := photo.NewRepo(sess)

	// -- initialize services ---
	userService := user.NewService(userRepo)
	storyService := story.NewService(storyRepo)
	commentService := comment.NewService(commentRepo)
	photoService := photo.NewService(uploadRepo)

	// URL/
	routes.PageRouter(app, userService, storyService)
	routes.AuthRouter(app, userService)

	// URL/api/
	api := app.Group("/api")
	routes.UserRouter(api, userService, storyService, commentService, photoService)
	routes.StoryRouter(api, userService, storyService, commentService, photoService)
	routes.CommentRouter(api, userService, storyService, commentService)
	routes.ImageRouter(api, photoService)

	// URL/api/admin/
	admin := api.Group("/admin")
	routes.AdminRouter(admin, userService, storyService)
}
