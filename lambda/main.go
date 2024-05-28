package main

import (
	"fmt"
	"lambda-func/app"
	"lambda-func/middleware"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Username string `json:"username"`
}

// Just take in a payload and do something with it
func handleRequest(event MyEvent) (string, error) {
	if event.Username == "" {
		return "", fmt.Errorf("Username can not be empty")
	}

	return fmt.Sprintf("%s\n", event.Username), nil
}

func ProtectedHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "User Unauthorized",
		StatusCode: http.StatusUnauthorized,
	}, nil
}

// This is basically out Backend
func main() {
	awsApp := app.NewApp()
	lambda.Start(func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch req.Path {
		case "/register":
			return awsApp.ApiHandler.RegisterUserHandler(req)
		case "/login":
			return awsApp.ApiHandler.LoginUser(req)
		case "/protected":
			return middleware.ValidateJWTMiddleware(ProtectedHandler)(req)
		default:
			return events.APIGatewayProxyResponse{
				Body:       "Not Found",
				StatusCode: http.StatusNotFound,
			}, nil
		}
	})
}
