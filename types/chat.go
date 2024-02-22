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
