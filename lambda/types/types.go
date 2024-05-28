package types

import (
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type NextFunction func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

func NewUser(registerUser RegisterUser) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), 10)
	if err != nil {
		return User{}, err
	}

	return User{
		Username:     registerUser.Username,
		PasswordHash: string(hashedPassword),
	}, nil
}

func ValidatePassword(hashedPassword string, plainTextPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainTextPassword))
	return err == nil
}

func CreateToken(user User) string {
	now := time.Now()
	validUntil := now.Add(time.Hour * 5).Unix()

	claims := jwtv5.MapClaims{
		"user":    user.Username,
		"expires": validUntil,
	}

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims, nil)
	secret := "storethesecretinawssecret"

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return fmt.Sprintf("%+v", err)
	}

	return tokenString

}
