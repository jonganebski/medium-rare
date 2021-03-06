package util

import (
	"home/jonganebski/github/medium-rare/config"
	"home/jonganebski/github/medium-rare/model"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GenerateCookie generates cookie storing jwt
func GenerateCookie(foundUser *model.User, exp time.Duration) (*fiber.Cookie, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = foundUser.ID
	claims["username"] = foundUser.Username
	claims["exp"] = time.Now().Add(exp)
	signedString, err := token.SignedString([]byte(config.Config("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "jwt"
	cookie.Value = signedString
	cookie.HTTPOnly = true
	cookie.Secure = config.Config("APP_ENV") == "PROD"
	cookie.Expires = time.Now().Add(exp)
	return cookie, nil
}

// GenerateCookieBeta generates cookie storing jwt
func GenerateCookieBeta(userOID *primitive.ObjectID, exp time.Duration) (*fiber.Cookie, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = userOID.Hex()
	claims["exp"] = time.Now().Add(exp)
	signedString, err := token.SignedString([]byte(config.Config("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "jwt"
	cookie.Value = signedString
	cookie.HTTPOnly = true
	cookie.Secure = config.Config("APP_ENV") == "PROD"
	cookie.Expires = time.Now().Add(exp)
	return cookie, nil
}
