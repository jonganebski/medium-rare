package router

import (
	"home/jonganebski/github/medium-rare/handler"
	"home/jonganebski/github/medium-rare/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up the routes
func SetupRoutes(app *fiber.App) {
	app.Get("/", handler.Home)
	app.Get("/new-story", middleware.Protected, handler.NewStory)
	app.Get("/read/:storyId", handler.ReadStory)
	app.Get("/edit-story/:storyId", handler.EditStory)
	app.Get("/followers/:userId", handler.SeeFollowers)

	app.Get("/signout", middleware.Protected, handler.Signout)
	app.Post("/signup", handler.CreateUser)
	app.Post("/signin", handler.Signin)

	me := app.Group("/me", middleware.Protected)
	me.Get("/bookmarks", handler.MyBookmarks)
	me.Get("/following", handler.SeeFollowings)
	me.Get("/settings", handler.SettingsPage)
	me.Get("/stories", handler.MyStories)

	publicAPI := app.Group("/api")
	publicAPI.Get("/blocks/:storyId", handler.ProvideStoryBlocks)
	publicAPI.Get("/comment/:storyId", handler.ProvideComments)

	privateAPI := app.Group("/api", middleware.Protected)

	privateAPI.Post("/bookmark/:storyId", handler.BookmarkStory)
	privateAPI.Post("/comment/:storyId", handler.AddComment)
	privateAPI.Post("/follow/:authorId", handler.Follow)
	privateAPI.Post("/like/:storyId/:plusMinus", handler.HandleLikeCount)
	privateAPI.Post("/photo/byfile", handler.UploadPhotoByFilename)
	privateAPI.Post("/unfollow/:authorId", handler.Unfollow)
	privateAPI.Post("/story", handler.AddStory)

	privateAPI.Patch("/story/:storyId", handler.UpdateStory)
	privateAPI.Patch("/user/username", handler.EditUsername)
	privateAPI.Patch("/user/bio", handler.EditBio)
	privateAPI.Patch("/user/avatar", handler.EditUserAvatar)
	privateAPI.Patch("/user/password", handler.EditPassword)

	privateAPI.Delete("/bookmark/:storyId", handler.DisBookmarkStory)
	privateAPI.Delete("/comment/:commentId", handler.DeleteComment)
	privateAPI.Delete("/photo", handler.DeletePhoto)
	privateAPI.Delete("/story/:storyId", handler.DeleteStory)
	privateAPI.Delete("/user", handler.DeleteUser)

	admin := app.Group("/admin")
	admin.Post("/pick/:storyId", handler.PickStory)
	admin.Post("/unpick/:storyId", handler.UnpickStory)
}
