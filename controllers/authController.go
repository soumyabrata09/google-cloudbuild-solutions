package controllers

import (
	"encoding/json"
	"net/http"

	"../models"
	u "../utils"
)

//These value methods will call the relevant functions from the models and eventually will call the CreateUser() and login() from our model
var CreateNewAccount = func(writer http.ResponseWriter, request *http.Request) {

	//getting the Account structure
	accountStruct := &models.Account{}
	//decoding the request body into the Account struct else throw an error otherwise
	accountCreationErr := json.NewDecoder(request.Body).Decode(accountStruct)
	if accountCreationErr != nil {
		u.Respond(writer, u.Message(false, "An Invalid Request Received"))
		return
	}
	//if no error occurred during the account creation phase , then we will proceed with calling our model method createUser
	//we're calling CreateUser() by concatinating it with the Account struct type reference as CreateUser() is a pointer receiver of type Struct
	response := accountStruct.CreateUser() // creating a new user and will pass this as an argument to the utility method Respond
	u.Respond(writer, response)
}

//These value methods will call the relevant functions from the models and eventually will call the Login() and login() from our model
var LetsAuthenticate = func(writer http.ResponseWriter, request *http.Request) {
	//getting the account here
	accountStruct := &models.Account{}
	// decoding the request body
	accountAuthAndLoginErr := json.NewDecoder(request.Body).Decode(accountStruct)
	if accountAuthAndLoginErr != nil {
		u.Respond(writer, u.Message(false, "Invalid Log in Credential"))
		return
	}
	//if no error occured during the login phase then simply route the the logged in page
	//Login() is just a simple function which is why we're only calling it directly
	responseBody := models.Login(accountStruct.Email, accountStruct.Password) // Here we're calling the Login() function from our model accounts.go
	u.Respond(writer, responseBody)
}
