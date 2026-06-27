package core_kafka



type UserRegisterEvent struct {
	ID  		string  `json:"user_id"`
	Username 	string  `json:"username"`
	PhoneNumber string  `json:"phone_number"`
}