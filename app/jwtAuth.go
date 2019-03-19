package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"../models"
	u "../utils"

	jwt "github.com/dgrijalva/jwt-go"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		notAuth := []string{"/api/user/new", "/api/user/login"} //List of endpoints that doesn't require auth
		requestPath := request.URL.Path                         //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(writer, request)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := request.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			response = u.Message(false, "Missing auth token")
			writer.WriteHeader(http.StatusForbidden)
			writer.Header().Add("Content-Type", "application/json")
			u.Respond(writer, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			response = u.Message(false, "Invalid/Malformed auth token")
			writer.WriteHeader(http.StatusForbidden)
			writer.Header().Add("Content-Type", "application/json")
			u.Respond(writer, response)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
		tokenModelData := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tokenModelData, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			response = u.Message(false, "Malformed authentication token")
			writer.WriteHeader(http.StatusForbidden)
			writer.Header().Add("Content-Type", "application/json")
			u.Respond(writer, response)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			response = u.Message(false, "Token is not valid.")
			writer.WriteHeader(http.StatusForbidden)
			writer.Header().Add("Content-Type", "application/json")
			u.Respond(writer, response)
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		fmt.Sprintf("User %s", tokenModelData.UserId) //Useful for monitoring
		ctx := context.WithValue(request.Context(), "user", tokenModelData.UserId)
		request = request.WithContext(ctx)
		next.ServeHTTP(writer, request) //proceed in the middleware chain!
	})
}
