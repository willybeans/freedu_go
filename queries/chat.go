package queries

import (
	"github.com/willybeans/freedu_go/database"
	"github.com/willybeans/freedu_go/logger"
)

type ChatRoomXref struct {
	ID      string `json:"id"`
	User_ID string `json:"user_id"`
	Chat_ID string `json:"chat_room_id"`
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
