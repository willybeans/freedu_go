// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websockets

import (
	"fmt"

	"github.com/willybeans/freedu_go/logger"
	"github.com/willybeans/freedu_go/pubsub"
	"github.com/willybeans/freedu_go/queries"
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
			for index, chat := range allChats {
				fmt.Println("index : ", index)
				fmt.Println("chat : ", chat)
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

			}
		case message := <-h.broadcast:
			for client := range h.clients {
				// println("test id, ", client.id)
				// add logic here that blocks messages from sending
				// if the specific client isnt subscribed to topic id

				// we need some switch logic here for
				// the different topics

				/*
				 - check for topic
				 - if not exist - make topic
				 then
				 subscribe user to topic
				*/

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
