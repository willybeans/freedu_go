// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websockets

import (
	"encoding/json"
	"fmt"

	"github.com/willybeans/freedu_go/logger"
	"github.com/willybeans/freedu_go/pubsub"
	"github.com/willybeans/freedu_go/queries"
	"github.com/willybeans/freedu_go/types"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	l := logger.Get()
	broker := pubsub.NewBroker()

	for { // this is infinite loop
		select { // allow you to wait on multiple channels
		//select BLOCKS untill one can run
		case client := <-h.register:
			println("REGISTER FIRED : ", client.id)

			h.clients[client] = true
			// SUBSCRIBE pubsub
			newSub := broker.AddSubscriber(client.id)
			// CHECK FOR TOPICS in db under this user
			allChats, err := queries.GetChatRoomsByUserID(client.id)
			//CreateChatPreviewsByUserID also needs to send the display messages
			if err != nil {
				l.Error().Err(err).Msg("Error Registering Socket Connection:")
			}
			for _, chat := range allChats {
				// fmt.Println("index : ", index)
				// fmt.Println("chat : ", chat)
				broker.Subscribe(newSub, chat.ChatRoom_ID)
			}

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				// UNSUBSCRIBE pubsub
				// broker.RemoveSubscriber() -- maybe not?
				// i think this is not needed since
				// we can keep them in memory for now ?
				// but we could use this for currently online state change

			}
		case message := <-h.broadcast:

			// Get struct from string.
			var webSocketMessage types.WebSocketMessage
			if err := json.Unmarshal(message, &webSocketMessage); err != nil {
				l.Error().Err(err).Msg("Error unmarshaling message:")
				continue
			}

			var responseBody types.ResponseBody
			switch webSocketMessage.Action {
			case "get_messages":
				fmt.Println("get_messages fired")
				// queries.GetMessagesByChatID(webSocketMessage.Content)
				//1- hit DB 2- spit back res obj
			case "post_message":
				l.Info().Msg("Post Message Websocket Fired")
				resObj, err := queries.NewMessageForUserInChat(types.NewMessage{ChatRoom_ID: webSocketMessage.ChatRoom_ID, User_ID: webSocketMessage.User_ID, Content: webSocketMessage.Content})
				if err != nil {
					l.Error().Err(err).Msg("Error NewMessage on Websocket Broadcast")
				}

				responseBody = types.ResponseBody{
					Message: resObj,
					Action:  webSocketMessage.Action,
				}

				b, err := json.Marshal(responseBody)
				if err != nil {
					l.Error().Err(err).Msg("Error unmarshaling message:")
				}

				message = b
			case "post_chat":
				fmt.Println("post_chat fired")
			default:
				fmt.Println("default fired")
			}

			for client := range h.clients {
				// println("test id, ", client.id)
				/*
				 - check for topic by message id
				 - if not exist - make topic
				 then
				 subscribe user to topic
				*/
				// if responseBody.Chatroom_ID

				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
					// DESTROY pubsub
				}
			}
		}
	}
}
