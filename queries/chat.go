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
	// use for building subscription
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

	var message types.Message
	query := database.DB().QueryRow(`	
	WITH new_m AS (
		INSERT INTO messages (chat_room_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING * )
	SELECT
		new_m.id,
		new_m.chat_room_id,
		new_m.user_id,
		new_m.content,
		new_m.sent_at,
		users.username
	FROM new_m JOIN users ON new_m.user_id = users.id`, newMessage.ChatRoom_ID, newMessage.User_ID, newMessage.Content)
	err := query.Scan(&message.ID, &message.ChatRoom_ID, &message.User_ID, &message.Content, &message.SentAt, &message.Username)
	if err != nil {
		l.Error().Err(err).Msg("Error NewMessage on Scan")
		return message, err
	}

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
	//please study this in more depth
	// ?is this query time bad compared to multiple joins?
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

func GetAllChatContentsByUserId(userId string) ([]types.UsersChats, error) {
	l := logger.Get()

	query := fmt.Sprintf(`
	WITH user_chat_rooms AS (
		SELECT
			ucr.chat_room_id,
			ucr.user_id
		FROM
			user_chatroom_xref ucr
		WHERE
			ucr.user_id = '%s'
	),
	messages_with_usernames AS (
		SELECT
			m.id AS message_id,
			m.chat_room_id,
			m.user_id,
			m.content,
			m.sent_at,
			u.username
		FROM
			messages m
			JOIN users u ON m.user_id = u.id
	)
	SELECT
		ucr.chat_room_id,
		cr.chat_name,
		ARRAY(SELECT username FROM users WHERE id IN (SELECT user_id FROM user_chatroom_xref WHERE chat_room_id = ucr.chat_room_id)) AS usernames,
		JSON_AGG(
			JSON_BUILD_OBJECT(
				'message_id', m.message_id,
				'content', m.content,
				'sent_at', m.sent_at,
				'username', m.username,
				'user_id', m.user_id
			) ORDER BY m.sent_at DESC
		) AS chat_messages
	FROM
		user_chat_rooms ucr
		JOIN chat_room cr ON ucr.chat_room_id = cr.id
		LEFT JOIN messages_with_usernames m ON ucr.chat_room_id = m.chat_room_id
	GROUP BY
		ucr.chat_room_id,
		cr.chat_name;
	`, userId)

	// Execute the query
	rows, err := database.DB().Query(query)
	if err != nil {
		l.Error().Err(err).Msg("Error GetAllChatsByUserId Query")
		return nil, err
	}
	defer rows.Close()

	// Process the results
	var chatrooms []types.UsersChats
	for rows.Next() {
		var chatroom types.UsersChats
		err := rows.Scan(&chatroom.Chatroom_ID, &chatroom.Chatroom_name, &chatroom.Usernames, &chatroom.ChatMessages)
		if err != nil {
			l.Error().Err(err).Msg("Error GetAllChatsByUserId Scan")
			return nil, err
		}
		chatrooms = append(chatrooms, chatroom)
	}

	return chatrooms, nil
}
