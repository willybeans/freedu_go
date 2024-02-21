package queries

import (
	"fmt"

	"github.com/lib/pq"
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

func NewXrefForChatID(ids types.IdsForNewXref) (types.ChatRoomXref, error) {
	l := logger.Get()

	// confirm user is allowed to write to this chat
	var xref types.ChatRoomXref
	// if userCanJoinChat(user) {
	query := database.DB().QueryRow("INSERT INTO user_chatroom_xref (chat_room_id, user_id) VALUES ($1, $2) RETURNING *", xref.ID, xref.ChatRoom_ID, xref.User_ID)
	err := query.Scan(&xref.ID, &xref.ChatRoom_ID, &xref.User_ID)
	if err != nil {
		l.Error().Err(err).Msg("Error NewMessage on Scan")
		return xref, err
	}
	// }
	return xref, nil

}

func CreateChatPreviewsByUserID(userId string) ([]types.PreviewMessage, error) {
	l := logger.Get()
	// join between chat id + members names (by xref) + recent message
	// ids := "'" + strings.Join(chatRoomIDs, "','") + "'"

	// Construct the query string
	query := fmt.Sprintf(`
	WITH recent_messages AS (
    SELECT DISTINCT ON (chat_room_id) *
    FROM messages
    ORDER BY chat_room_id, sent_at DESC
	)
	SELECT
    u.username,
    rm.content,
    rm.sent_at,
    ARRAY(SELECT username FROM users WHERE id IN (SELECT user_id FROM user_chatroom_xref WHERE chat_room_id = rm.chat_room_id)) AS usernames
	FROM
    recent_messages rm
	JOIN
    users u ON rm.user_id = u.id
	WHERE
    rm.chat_room_id IN (SELECT chat_room_id FROM user_chatroom_xref WHERE user_id = '%s');
	`, userId)

	// Execute the query
	rows, err := database.DB().Query(query)
	if err != nil {
		l.Error().Err(err).Msg("Error JoinChatUserMessagePreview Query")
		return nil, err
	}
	defer rows.Close()

	// Process the results
	var messages []types.PreviewMessage
	for rows.Next() {
		var message types.PreviewMessage
		err := rows.Scan(&message.Username, &message.Content, &message.SentAt, pq.Array(&message.UserNames))
		if err != nil {
			l.Error().Err(err).Msg("Error JoinChatUserMessagePreview Scan")
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil

}
