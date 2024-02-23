package types

type WebSocketMessage struct {
	Action      string `json:"action"`
	ChatRoom_ID string `json:"chat_room_id"`
	User_ID     string `json:"user_id"`
	Content     string `json:"content"`
}

type ResponseBody struct {
	Message
	Action string `json:"action"`
}
