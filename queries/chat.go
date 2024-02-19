package queries

import (
	"github.com/willybeans/freedu_go/database"
	"github.com/willybeans/freedu_go/logger"
	"github.com/willybeans/freedu_go/types"
)

func GetMessagesByChatID(chatRoomID string) ([]types.Message, error) {
	l := logger.Get()

	rows, err := database.DB().Query("SELECT * FROM messages WHERE chat_room_id = $1 ORDER BY \"sent_at\"", chatRoomID)
	if err != nil {
		l.Error().Err(err).Msg("Error GetMessagesByChatID on Query")
		return nil, err
	}
	defer rows.Close()

	Messages := make([]types.Message, 0)
	for rows.Next() {
		var message types.Message
		err := rows.Scan(&message.ID, &message.ChatRoom_ID, &message.User_ID, &message.Content, &message.SentAt)
		if err != nil {
			l.Error().Err(err).Msg("Error GetMessagesByChatID on Scan")
			return nil, err
		}
		Messages = append(Messages, message)
	}

	return Messages, nil
}

func GetXRefsByChatID(chatId string) ([]types.ChatRoomXref, error) {
	l := logger.Get()

	rows, err := database.DB().Query("SELECT * FROM user_chatroom_xref WHERE chat_room_id = $1", chatId)
	if err != nil {
		l.Error().Err(err).Msg("Error GetXRefsByChatID on Query")
		return nil, err
	}
	defer rows.Close()

	chatRoomXrefList := make([]types.ChatRoomXref, 0)
	for rows.Next() {
		var chatRoomXref types.ChatRoomXref
		err := rows.Scan(&chatRoomXref.ID, &chatRoomXref.User_ID, &chatRoomXref.ChatRoom_ID)
		if err != nil {
			l.Error().Err(err).Msg("Error GetXRefsByChatID on Scan")
			return nil, err
		}
		chatRoomXrefList = append(chatRoomXrefList, chatRoomXref)
	}
	return chatRoomXrefList, nil
}

func GetChatRoomsByUserID(userId string) ([]types.ChatRoomXref, error) {
	l := logger.Get()

	rows, err := database.DB().Query("SELECT * FROM user_chatroom_xref WHERE user_id = $1", userId)
	if err != nil {
		l.Error().Err(err).Msg("Error GetChatRoomsByUserID on Query")
		return nil, err
	}
	defer rows.Close()

	chatRoomXrefList := make([]types.ChatRoomXref, 0)
	for rows.Next() {
		var chatRoomXref types.ChatRoomXref
		err := rows.Scan(&chatRoomXref.ID, &chatRoomXref.User_ID, &chatRoomXref.ChatRoom_ID)
		if err != nil {
			l.Error().Err(err).Msg("Error GetChatRoomsByUserID on Scan")
			return nil, err
		}
		chatRoomXrefList = append(chatRoomXrefList, chatRoomXref)
	}
	return chatRoomXrefList, nil
}

func NewMessageForUserInChat(newMessage types.NewMessage) (types.Message, error) {
	l := logger.Get()

	// confirm user is allowed to write to this chat
	var message types.Message
	// if userCanWriteToChat(newMessage) {
	query := database.DB().QueryRow("INSERT INTO messages (chat_room_id, user_id, content) VALUES ($1, $2, $3) RETURNING *", newMessage.ChatRoom_ID, newMessage.User_ID, newMessage.Content)
	err := query.Scan(&message.ID, &message.ChatRoom_ID, &message.User_ID, &message.Content, &message.SentAt)
	if err != nil {
		l.Error().Err(err).Msg("Error NewMessage on Scan")
		return message, err
	}
	// }
	return message, nil

}
