package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func RenderHome(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "views/profile.html")
}

func RenderLogin(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "views/login.html")
}

func RenderRegister(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "views/register.html")
}

// SignInUser Used for Signing In the Users
func SignInUser(response http.ResponseWriter, request *http.Request) {
	var loginRequest LoginParams
	var result UserDetails
	var errorResponse = ErrorResponse{
		Code: http.StatusInternalServerError, Message: "Incorrect Credentials",
	}

	decoder := json.NewDecoder(request.Body)
	decoderErr := decoder.Decode(&loginRequest)
	defer request.Body.Close()

	if decoderErr != nil {
		returnErrorResponse(response, request, errorResponse)
	} else {
		errorResponse.Code = http.StatusBadRequest
		if loginRequest.Email == "" {
			errorResponse.Message = "Email can't be empty"
			returnErrorResponse(response, request, errorResponse)
		} else if loginRequest.Password == "" {
			errorResponse.Message = "Password can't be empty"
			returnErrorResponse(response, request, errorResponse)
		} else {

			collection := Client.Database("test").Collection("users")

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			var err = collection.FindOne(ctx, bson.M{
				"email":    loginRequest.Email,
				"password": loginRequest.Password,
			}).Decode(&result)

			defer cancel()

			if err != nil {
				returnErrorResponse(response, request, errorResponse)
			} else {
				tokenString, _ := CreateJWT(loginRequest.Email)

				if tokenString == "" {
					returnErrorResponse(response, request, errorResponse)
				}

				var successResponse = SuccessResponse{
					Code:    http.StatusOK,
					Message: "Login Sucessfull !",
					Response: SuccessfulLoginResponse{
						AuthToken: tokenString,
						Email:     loginRequest.Email,
					},
				}

				successJSONResponse, jsonError := json.Marshal(successResponse)

				if jsonError != nil {
					returnErrorResponse(response, request, errorResponse)
				}
				response.Header().Set("Content-Type", "application/json")
				response.Write(successJSONResponse)
			}
		}
	}
}

// SignUpUser Used for Signing up the Users
func SignUpUser(response http.ResponseWriter, request *http.Request) {
	var registationRequest RegistationParams
	var errorResponse = ErrorResponse{
		Code: http.StatusInternalServerError, Message: "Security breach !!",
	}

	decoder := json.NewDecoder(request.Body)
	decoderErr := decoder.Decode(&registationRequest)
	defer request.Body.Close()

	if decoderErr != nil {
		returnErrorResponse(response, request, errorResponse)
	} else {
		errorResponse.Code = http.StatusBadRequest
		if registationRequest.Name == "" {
			errorResponse.Message = "Username can't be empty"
			returnErrorResponse(response, request, errorResponse)
		} else if registationRequest.Email == "" {
			errorResponse.Message = "Email can't be empty"
			returnErrorResponse(response, request, errorResponse)
		} else if registationRequest.Password == "" || len(registationRequest.Password) < 8 {

			errorResponse.Message = "Password can't be empty"
			returnErrorResponse(response, request, errorResponse)
		} else if registationRequest.Conpass == "" || len(registationRequest.Conpass) < 8 || registationRequest.Password != registationRequest.Conpass {
			errorResponse.Message = "Error confirming password. Make sure its length is greater then 8 and its same as your password"
			returnErrorResponse(response, request, errorResponse)
		} else if registationRequest.Phone == "" {
			errorResponse.Message = "Phone number can't be empty"
			returnErrorResponse(response, request, errorResponse)
		} else {

			tokenString, _ := CreateJWT(registationRequest.Email)

			if tokenString == "" {
				returnErrorResponse(response, request, errorResponse)
			}
			var registrationResponse = SuccessfulLoginResponse{
				AuthToken: tokenString,
				Email:     registationRequest.Email,
			}

			collection := Client.Database("test").Collection("old")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			_, databaseErr := collection.InsertOne(ctx, bson.M{
				"name":     registationRequest.Name,
				"email":    registationRequest.Email,
				"password": registationRequest.Password,
				"conpass":  registationRequest.Conpass,
				"phone":    registationRequest.Phone,
				"invite":   registationRequest.Invite,
			})

			if databaseErr != nil {
				errorResponse.Code = http.StatusBadRequest
				fmt.Println("*********************DUPLICATE RECORD ERROR********************")
				var result UserDetails
				var emailDuplicationErr = collection.FindOne(ctx, bson.M{
					"email": registationRequest.Email,
				}).Decode(&result)

				if emailDuplicationErr == nil {
					errorResponse.Code = http.StatusBadRequest
					errorResponse.Message = "Email already exists"
					returnErrorResponse(response, request, errorResponse)
					log.Println("Email already exits")
				}
				var phoneDuplicationErr = collection.FindOne(ctx, bson.M{
					"phone": registationRequest.Phone,
				}).Decode(&result)
				if phoneDuplicationErr == nil {
					log.Println("Phone number already exits")
				}
				var nameDuplicationErr = collection.FindOne(ctx, bson.M{
					"name": registationRequest.Name,
				}).Decode(&result)
				if nameDuplicationErr == nil {
					log.Println("Username already exits")
				}
				var inviteDuplicationErr = collection.FindOne(ctx, bson.M{
					"invite": registationRequest.Invite,
				}).Decode(&result)
				if inviteDuplicationErr == nil {
					log.Println("Invitation code already used")
				}

			}

			defer cancel()

			var successResponse = SuccessResponse{
				Code:     http.StatusOK,
				Message:  "Registration sucessfull !",
				Response: registrationResponse,
			}

			successJSONResponse, jsonError := json.Marshal(successResponse)

			if jsonError != nil {
				returnErrorResponse(response, request, errorResponse)
			}
			response.Header().Set("Content-Type", "application/json")
			response.WriteHeader(successResponse.Code)
			response.Write(successJSONResponse)
		}
	}

}

// GetUserDetails Used for getting the user details using user token
func GetUserDetails(response http.ResponseWriter, request *http.Request) {
	var result UserDetails
	var errorResponse = ErrorResponse{
		Code: http.StatusInternalServerError, Message: "Security breach!",
	}
	bearerToken := request.Header.Get("Authorization")
	var authorizationToken = strings.Split(bearerToken, " ")[1]

	email, _ := VerifyToken(authorizationToken)
	if email == "" {
		returnErrorResponse(response, request, errorResponse)
	} else {
		collection := Client.Database("test").Collection("users")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var err = collection.FindOne(ctx, bson.M{
			"email": email,
		}).Decode(&result)

		defer cancel()

		if err != nil {
			returnErrorResponse(response, request, errorResponse)
		} else {
			var successResponse = SuccessResponse{
				Code:     http.StatusOK,
				Message:  "You are logged in successfully",
				Response: result.Name,
			}

			successJSONResponse, jsonError := json.Marshal(successResponse)

			if jsonError != nil {
				returnErrorResponse(response, request, errorResponse)
			}
			response.Header().Set("Content-Type", "application/json")
			response.Write(successJSONResponse)
		}
	}
}

func returnErrorResponse(response http.ResponseWriter, request *http.Request, errorMesage ErrorResponse) {
	httpResponse := &ErrorResponse{Code: errorMesage.Code, Message: errorMesage.Message}
	jsonResponse, err := json.Marshal(httpResponse)
	if err != nil {
		panic(err)
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(errorMesage.Code)
	response.Write(jsonResponse)
}
