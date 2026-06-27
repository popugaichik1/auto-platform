package auth_transport_http

import "github.com/google/uuid"


type UserDTOResponse struct {
	ID 				uuid.UUID  	`json:"id"`
	Version 		int   		`json:"version"`
	Username 		string 		`json:"username"` 
	PhoneNumber 	string 		`json:"phone_number"`
}


type CreateUserRequest struct {
	Username 		string `json:"username" binding:"required,min=3,max=100"` 
	PhoneNumber 	string `json:"phone_number" binding:"omitempty,min=10,max=15,startswith=+"`
	Password 		string `json:"password"     binding:"required"`
}


type LoginRequest struct {
	PhoneNumber 	string `json:"phone_number" binding:"required,min=10,max=15,startswith=+"`
	Password 		string `json:"password" binding:"required,min=8"` 
}


type LoginResponse struct {
	AccessToken  		string  `json:"access_token"`
	RefreshToken 		string  `json:"refresh_token"`
}


type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	Token string `json:"access_token"`
}
