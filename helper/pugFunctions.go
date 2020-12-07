package helper

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// IsPublishButton determines publish button's existance
func IsPublishButton(path string) bool {
	if strings.Contains(path, "new-story") {
		return true
	}
	if strings.Contains(path, "edit-story") {
		return true
	}
	return false
}

// GetPublishBtnText determines publish button's text
func GetPublishBtnText(path string) string {
	if strings.Contains(path, "new-story") {
		return "Publish"
	}
	if strings.Contains(path, "edit-story") {
		return "Save and publish"
	}
	return ""
}

// IsMyStory tells target story's author is th user or not
func IsMyStory(authorID, userID string) bool {
	if authorID == userID {
		return true
	}
	return false
}

// GetStoryPostDate translates unix based post date
func GetStoryPostDate(createdAt int64) string {
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
}

// GrindBody cuts and refine long body text
func GrindBody(body string, targetLen int) string {
	if targetLen < len(body) {
		return body[:targetLen] + "..."
	}
	return body
}

// GetSliceLen returns formatted count in string
func GetSliceLen(slice []primitive.ObjectID) string {
	len := len(slice)
	p := message.NewPrinter(language.English)
	formatted := p.Sprintf("%v", len)
	return formatted
}
