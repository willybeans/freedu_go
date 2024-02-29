package types

import "encoding/json"

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
	Username    string `json:"username"`
}

type ChatRoomXref struct {
	ID          string `json:"id"`
	User_ID     string `json:"user_id"`
	ChatRoom_ID string `json:"chat_room_id"`
}

type IdsForNewXref struct {
	User_ID     string `json:"user_id"`
	ChatRoom_ID string `json:"chat_room_id"`
}

// type ChatPreview struct {
// 	Members     []string       `json:"members"`
// 	ChatRoom_ID string         `json:"chat_room_id"`
// 	Preview     PreviewMessage `json:preview`
// }

type PreviewMessage struct {
	Username  string   `json:"username"`
	Content   string   `json:"content"`
	SentAt    string   `json:"sent_at"`
	UserNames []string `json:"usernames"`
}

type UsersChats struct {
	Chatroom_ID   string          `json:"chat_room_id"`
	Chatroom_name string          `json:"chat_name"`
	ChatMessages  json.RawMessage `json:"chat_messages"`
}

// consider combing with messages
type ChatMessage struct {
	ID             string `json:"id"`
	Content        string `json:"content"`
	SentAt         string `json:"sent_at"`
	SenderUsername string `json:"sender_username"`
}
