package router

import (
	"home/jonganebski/github/medium-rare/handler"
	"home/jonganebski/github/medium-rare/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up the routes
func SetupRoutes(app *fiber.App) {
	app.Get("/", handler.Home)                                               // refactored
	app.Get("/new-story", middleware.Protected, handler.NewStory)            // refactored
	app.Get("/read/:storyId", handler.ReadStory)                             // refactored
	app.Get("/edit-story/:storyId", middleware.Protected, handler.EditStory) // refactored
	app.Get("/followers/:userId", handler.SeeFollowers)                      // refactored.. kind of
	app.Get("/user-home/:userId", handler.UserHome)                          // refactored

	app.Get("/signout", middleware.Protected, handler.Signout) // refactored
	app.Post("/signup", handler.CreateUser)                    // refactored
	app.Post("/signin", handler.Signin)                        // refactored

	me := app.Group("/me", middleware.Protected)
	me.Get("/bookmarks", handler.MyBookmarks)   // refactored
	me.Get("/following", handler.SeeFollowings) // refactored
	me.Get("/settings", handler.SettingsPage)   // refactored
	me.Get("/stories", handler.MyStories)       // refactored

	publicAPI := app.Group("/api")
	publicAPI.Get("/blocks/:storyId", handler.ProvideStoryBlocks) // refactored
	publicAPI.Get("/comment/:storyId", handler.ProvideComments)   // refactored

	privateAPI := app.Group("/api", middleware.APIGuard)

	privateAPI.Post("/bookmark/:storyId", handler.BookmarkStory)          // refactored
	privateAPI.Post("/comment/:storyId", handler.AddComment)              // refactored
	privateAPI.Post("/follow/:authorId", handler.Follow)                  // refactored
	privateAPI.Post("/like/:storyId/:plusMinus", handler.HandleLikeCount) // refactored
	privateAPI.Post("/photo/byfile", handler.UploadPhotoByFilename)       // refactored
	privateAPI.Post("/unfollow/:authorId", handler.Unfollow)              // refactored
	privateAPI.Post("/story", handler.AddStory)                           // refactored

	privateAPI.Patch("/story/:storyId", handler.UpdateStory) // refactored
	privateAPI.Patch("/user/username", handler.EditUsername) // refactored
	privateAPI.Patch("/user/bio", handler.EditBio)           // refactored
	privateAPI.Patch("/user/avatar", handler.EditUserAvatar) // refactored
	privateAPI.Patch("/user/password", handler.EditPassword) // refactored

	privateAPI.Delete("/bookmark/:storyId", handler.DisBookmarkStory) // refactored
	privateAPI.Delete("/comment/:commentId", handler.DeleteComment)   // refactored
	privateAPI.Delete("/photo", handler.DeletePhoto)
	privateAPI.Delete("/story/:storyId", handler.DeleteStory) // refactored
	privateAPI.Delete("/user", handler.DeleteUser)            // refactored

	admin := app.Group("/admin")
	admin.Post("/pick/:storyId", handler.PickStory)
	admin.Post("/unpick/:storyId", handler.UnpickStory)
}
