package router

import (
	"home/jonganebski/github/medium-rare/config"
	"home/jonganebski/github/medium-rare/database"
	"home/jonganebski/github/medium-rare/package/comment"
	"home/jonganebski/github/medium-rare/package/story"
	"home/jonganebski/github/medium-rare/package/user"
	"home/jonganebski/github/medium-rare/routes"

	"github.com/gofiber/fiber/v2"
)

var mg = &database.Mongo

// SetupRoutes sets up the routes
func SetupRoutes(app *fiber.App) {
	// --- mongodb collections ---
	userCollection := mg.Db.Collection(config.Config("COLLECTION_USER"))
	storyCollection := mg.Db.Collection(config.Config("COLLECTION_STORY"))
	commentCollection := mg.Db.Collection(config.Config("COLLECTION_COMMENT"))
	// --- initialized repositories with collections ---
	userRepo := user.NewRepo(userCollection)
	storyRepo := story.NewRepo(storyCollection)
	commentRepo := comment.NewRepo(commentCollection)
	// -- initialize services with repositories ---
	userService := user.NewService(userRepo)
	storyService := story.NewService(storyRepo)
	commentService := comment.NewService(commentRepo)

	// URL/
	routes.PageRouter(app, userService, storyService)
	routes.AuthRouter(app, userService)

	// URL/api/
	api := app.Group("/api")
	routes.UserRouter(api, userService, storyService, commentService)
	routes.StoryRouter(api, userService, storyService, commentService)
	routes.CommentRouter(api, userService, storyService, commentService)
	routes.ImageRouter(api)

	// URL/api/admin/
	admin := api.Group("/admin")
	routes.AdminRouter(admin, userService, storyService)
}
