package main

import jwt "github.com/dgrijalva/jwt-go"

type ErrorResponse struct {
	Code    int
	Message string
}

type SuccessResponse struct {
	Code     int
	Message  string
	Response interface{}
}

type Claims struct {
	Email string
	jwt.StandardClaims
}

type RegistationParams struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Conpass  string `json:"conpass"`
	Phone    string `json:"phone"`
	Invite   string `json:"invite"`
}

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SuccessfulLoginResponse struct {
	Email     string
	AuthToken string
}

type UserDetails struct {
	Name     string
	Email    string
	Password string
}
