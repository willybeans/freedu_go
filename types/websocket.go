package types

type WebSocketMessage struct {
	Action      string `json:"action"`
	ChatRoom_ID string `json:"chat_room_id"`
	User_ID     string `json:"user_id"`
	Content     string `json:"content"`
}

type ResponseBody struct {
	SubscriptionID string      `json:"subscription_id"`
	Content        interface{} `json:"content"`
	// Message        Message      `json:"message"`
	// UserChats      []UsersChats `json:"user_chats"`
	Action string `json:"action"`
}
