package queries

import (
	"github.com/willybeans/freedu_go/database"
	"github.com/willybeans/freedu_go/logger"
)

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

func GetMessagesByChatID(chatRoomID string) ([]Message, error) {
	l := logger.Get()

	rows, err := database.DB().Query("SELECT * FROM messages WHERE chat_room_id = $1 ORDER BY \"sent_at\"", chatRoomID)
	if err != nil {
		l.Error().Err(err).Msg("Error GetMessagesByChatID on Query")
		return nil, err
	}
	defer rows.Close()

	Messages := make([]Message, 0)
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.ID, &message.ChatRoom_ID, &message.User_ID, &message.Content, &message.SentAt)
		if err != nil {
			l.Error().Err(err).Msg("Error GetMessagesByChatID on Scan")
			return nil, err
		}
		Messages = append(Messages, message)
	}

	return Messages, nil
}

func GetChatRoomsByUserID(userId string) ([]ChatRoomXref, error) {
	l := logger.Get()

	rows, err := database.DB().Query("SELECT * FROM user_chatroom_xref WHERE user_id = $1", userId)
	if err != nil {
		l.Error().Err(err).Msg("Error GetChatRoomsByUserID on Query")
		return nil, err
	}
	defer rows.Close()

	chatRoomXrefList := make([]ChatRoomXref, 0)
	for rows.Next() {
		var chatRoomXref ChatRoomXref
		err := rows.Scan(&chatRoomXref.ID, &chatRoomXref.User_ID, &chatRoomXref.Chat_ID)
		if err != nil {
			l.Error().Err(err).Msg("Error GetChatRoomsByUserID on Scan")
			return nil, err
		}
		chatRoomXrefList = append(chatRoomXrefList, chatRoomXref)
	}
	return chatRoomXrefList, nil
}
