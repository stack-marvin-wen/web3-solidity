package utils

import (
	"oauth/config"
	"oauth/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(u models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  u.ID,
		"email":    u.Email,
		"username": u.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Config.JWT_SECRET))
}

func ParseJWT(token string) (models.User, error) {
	var user models.User
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.Config.JWT_SECRET), nil
	})
	if err != nil {
		return user, err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		user.ID = uint(claims["user_id"].(float64))
		user.Email = claims["email"].(string)
		user.Username = claims["username"].(string)
		return user, nil
	}
	return user, jwt.ErrInvalidKey
}
