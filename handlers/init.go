package handlers

import (
	"net/http"

	"github.com/hsputra/API-CGPT/types"

	"github.com/gorilla/websocket"
)

var (
	//  websocket upgrader
	upgrader = websocket.Upgrader{
		// check origin of the request
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// declare connection pool variable with reference to NewConnectionPool types
var connectionPool = types.NewConnectionPool()

// declare connection pool variable with reference to NewConversationPool types
var conversationPool = types.NewConversationPool()
