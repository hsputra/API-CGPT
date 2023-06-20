package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hsputra/API-CGPT/types"
	"github.com/hsputra/API-CGPT/utils"
)

// Client routes to handle client connections
func Client_register(c *gin.Context) {

	// make websocket connection
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	// generate connection id
	id := utils.GenerateId()
	// create err variable and send connection id with types.Message as the parameter
	err = ws.WriteJSON(types.Message{
		Id:      id,
		Message: "connection id",
	})
	if err != nil {
		return
	}

	// wait for client to connection id by using for loop
	for {
		// read message from client
		var msg types.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			return
		}

		// check if message id is connection id
		if msg.Id == id {
			break
		} else {
			// reconnect if message id is not connection id
			// check if connection id is already in connection pool
			connection, ok := connectionPool.Get(msg.Id)
			if ok {
				// close the old connection
				connection.Ws.Close()
				// delete the old connection from connection pool
				connectionPool.Delete(msg.Id)
			}
			// assign message id to connection id
			id = msg.Id
			break
		}
	}

	// add connection to connection pool
	connection := &types.Connection{
		Id: id,

		// set websocket connection
		Ws: ws,

		// set last message time to beginning of time
		LastMessageTime: time.Time{},
		Heartbeat:       time.Now(),
	}
	connectionPool.Set(connection)
	// Debug
	println("New connection", connection.Id)
}
