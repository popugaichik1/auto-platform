package transport_http

import "github.com/google/uuid"


type SaveUserRequest struct {
	ID 				uuid.UUID  	`json:"id" validate:"required"`
	Username 		string     	`json:"username" validate:"required,min=1,max=125"`
	PhoneNumber 	string	   	`json:"phone_number" validate:"required,min=10,max-15,startswith=+"`
}

type UserDTO struct {
	ID 				uuid.UUID  	`json:"id"`
	Version         int   		`json:"version"`
	Username 		string    	`json:"username"`
	PhoneNumber 	string		`json:"phone_number"`
}