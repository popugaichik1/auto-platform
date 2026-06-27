package core_kafka

type Message struct {
	Topic 		[]byte
	Key 		[]byte
	Payload 	any
}


func NewMessage(topic, key []byte, payload any) Message {
	return Message{
		Topic: 		topic,
		Key: 		key,
		Payload: 	payload,
	}
}
