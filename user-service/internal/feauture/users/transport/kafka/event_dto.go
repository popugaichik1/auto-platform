package transport_kafka

import "time"


type UserRegisterEvent struct {
	ID          string `json:"user_id"`
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
}

// DeadLetterEvent оборачивает сообщение, которое не удалось обработать
// без шанса на успех при повторной попытке (битый JSON, невалидные данные).
// Публикуется в TopicUserRegisteredDLQ для последующего ручного разбора.
type DeadLetterEvent struct {
	OriginalTopic     string    `json:"original_topic"`
	OriginalPartition int32     `json:"original_partition"`
	OriginalOffset    int64     `json:"original_offset"`
	Error             string    `json:"error"`
	FailedAt          time.Time `json:"failed_at"`
	Payload           string    `json:"payload"`
}

