package auth

import (
	"fmt"
	"os"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // Expires in 7 days
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func GetUserIDFromToken(c *gin.Context) (uint, error) {
	tokenString, err := c.Cookie("jwt")
	if err != nil {
		return 0, fmt.Errorf("could not retrieve token from cookie")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub := claims["sub"].(float64)
		userID := uint(sub)
		return userID, nil
	}
	
	return 0, fmt.Errorf("invalid token")
}
