package JWT

import (
	"IosifSuzuki/sharingToMe/internal/configuration"
	"IosifSuzuki/sharingToMe/internal/models"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

func GenerateToken(consumer models.Consumer) (string, error) {
	var tokenClaim = models.Token{
		Username:    consumer.Username,
		PhoneNumber: consumer.PhoneNumber,
		Role: "user",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	}
	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaim)
	return token.SignedString([]byte(configuration.Configuration.AppInfo.SecretKey))
}

func ValidateToken(bearerToken string) (*jwt.Token, error) {
	tokenString := strings.Split(bearerToken, " ")[1]
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
			return nil, errors.New("Expired token")
		}
		return []byte(configuration.Configuration.AppInfo.SecretKey), nil
	})
}
