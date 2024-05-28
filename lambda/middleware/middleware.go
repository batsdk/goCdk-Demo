package middleware

import (
	"fmt"
	"lambda-func/types"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// Extract req headers
// Extract req claims
// Then validate

func ValidateJWTMiddleware(next types.NextFunction) types.NextFunction {
	return func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

		// Extract Headers
		tokenString := extractTokenFromHeaders(req.Headers)
		if tokenString == "" {
			return events.APIGatewayProxyResponse{
				Body:       "Auth Headers are Missing",
				StatusCode: http.StatusUnauthorized,
			}, fmt.Errorf("Auth Headers are missing")
		}

		// Parse -> Token Claims
		claims, err := parseToken(tokenString)
		if err != nil {
			return events.APIGatewayProxyResponse{
				Body:       "User Unauthorized",
				StatusCode: http.StatusUnauthorized,
			}, err
		}

		expires := int64(claims["expires"].(float64))
		if time.Now().Unix() > expires {
			return events.APIGatewayProxyResponse{
				Body:       "Token Expired",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}

		return next(req)

	}
}

func extractTokenFromHeaders(headers map[string]string) string {
	authHeader, ok := headers["Authorization"]

	if !ok {
		return ""
	}

	splitToken := strings.Split(authHeader, "Bearer ")

	if len(splitToken) != 2 {
		return ""
	}

	return splitToken[2]
}

func parseToken(tokenString string) (jwtv5.MapClaims, error) {
	secret := "storethesecretinawssecret"

	token, err := jwtv5.Parse(tokenString, func(token *jwtv5.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("StatusUnauthorized")
	}

	if !token.Valid {
		return nil, fmt.Errorf("Token is not valid")
	}

	claims, ok := token.Claims.(jwtv5.MapClaims)
	if !ok {
		return nil, fmt.Errorf("claims of unauthorized type")
	}

	return claims, nil
}
