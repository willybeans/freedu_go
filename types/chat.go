package types

type NewMessage struct {
	ChatRoom_ID string `json:"chat_room_id"`
	User_ID     string `json:"user_id"`
	Content     string `json:"content"`
}

type Message struct {
	ID          string `json:"id"`
	ChatRoom_ID string `json:"chat_room_id"`
	User_ID     string `json:"user_id"`
	Content     string `json:"content"`
	SentAt      string `json:"sent_at"`
}

type ChatRoomXref struct {
	ID      string `json:"id"`
	User_ID string `json:"user_id"`
	Chat_ID string `json:"chat_room_id"`
}
