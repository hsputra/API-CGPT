package handlers

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hsputra/API-CGPT/types"
	"github.com/hsputra/API-CGPT/utils"
)

// API Routes
func API_ask(c *gin.Context) {

	//  get request with ChatGptRequest as type
	var request types.ChatGptRequest
	err := c.BindJSON(&request)
	if err != nil {
		// error 400 Invalid request body
		c.JSON(400, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// check if authorization is in the header
	if c.Request.Header["Authorization"] == nil {
		// error 401 Unauthorized
		c.JSON(401, gin.H{
			"error": "API key not provided",
		})
		return
	}

	//  check if API key is valid
	verified, err := utils.VerifyToken(c.Request.Header["Authorization"][0])
	if err != nil {
		// error 500 Internal server error
		c.JSON(500, gin.H{
			"error": "Failed to verify API key",
		})
		return
	}
	if !verified {
		// error 401 Unauthorized
		c.JSON(401, gin.H{
			"error": "Invalid API key",
		})
		return
	}

	// if message id is not provided, generate message id
	if request.MessageId == "" {
		request.MessageId = utils.GenerateId()
	}

	// if parent id is not set, generate a new one
	if request.ParentId == "" {
		request.ParentId = utils.GenerateId()
	}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		// error 500 Internal server error
		c.JSON(500, gin.H{
			"error": "Failed to convert request to json",
		})
		return
	}

	var connection *types.Connection
	// check conversation id
	connectionPool.Mu.RLock()
	// check number of connections
	if len(connectionPool.Connections) == 0 {
		// error 503 Internal server error
		c.JSON(503, gin.H{
			"error": "No available clientszzz",
		})
		return
	}
	connectionPool.Mu.RUnlock()
	if request.ConversationId == "" {
		// retry 3 times
		var succeeded bool = false
		for i := 0; i < 3; i++ {
			// find connection with the lowest load and where heartbeat is after last message time
			connectionPool.Mu.RLock()
			// for loop to find connections
			for _, conn := range connectionPool.Connections {
				// check if connection is nil or last message time is after heartbeat
				if connection == nil || conn.LastMessageTime.Before(connection.LastMessageTime) {
					// check if connection heartbeat is after last message time
					if conn.Heartbeat.After(conn.LastMessageTime) {
						connection = conn
					}
				}
			}
			connectionPool.Mu.RUnlock()
			// check if connection was found
			if connection != nil {
				// return 503 for no available clients
				c.JSON(503, gin.H{
					"error": "No available clients",
				})
				return
			}

			// ping before sending request
			var pingSucceeded bool = ping(connection.Id)
			if !pingSucceeded {
				// ping failed, try again
				connectionPool.Delete(connection.Id)
				succeeded = false
				connection = nil
				continue
			} else {
				succeeded = true
				break
			}
		}
		if !succeeded {
			// delete connection
			// return 503 for failed ping
			c.JSON(503, gin.H{
				"error": "Ping failed",
			})
			return
		}
	} else {
		//  check if conversation exists
		conversation, ok := conversationPool.Get(request.ConversationId)
		if !ok {
			// error 500
			c.JSON(500, gin.H{
				"error": "Conversation does not exist",
			})
			return
		} else {
			// get connectionId of the conversation
			connectionId := conversation.ConnectionId
			// check if the connection exists
			connection, ok := connectionPool.Get(connectionId)
			if !ok {
				// error 500
				c.JSON(500, gin.H{
					"error": "Connection no longer exist",
				})
				return
			}
			//  ping before sending request
			if !ping(connection.Id) {

				// return 503 for failed ping
				c.JSON(503, gin.H{
					"error": "Ping failed",
				})
				return
			}
		}

		message := types.Message{
			Id:      utils.GenerateId(),
			Message: "API-CGPT Request",
			// convert request to json
			Data: string(jsonRequest),
		}
		err := connection.Ws.WriteJSON(message)
		if err != nil {
			// error 500
			c.JSON(500, gin.H{
				"error": "Failed to send request",
			})
			// delete connection
			connectionPool.Delete(connection.Id)
			return
		}

		// set last message time
		connection.LastMessageTime = time.Now()
		// wait for response with timeout
		for {
			// read message
			var receive types.Message
			connection.Ws.SetReadDeadline(time.Now().Add(120 * time.Second))
			err = connection.Ws.ReadJSON(&receive)
			if err != nil {
				// error 500
				c.JSON(500, gin.H{
					"error": "Failed to read response",
				})
				// delete connection
				connectionPool.Delete(connection.Id)
				return
			}
			// check if message is the response
			if receive.Id == message.Id {
				// convert response to ChatGptResponse
				var response types.ChatGptResponse
				err = json.Unmarshal([]byte(receive.Data), &response)
				if err != nil {
					// error 500
					c.JSON(500, gin.H{
						"error":    "Failed to convert response to CGPT response",
						"response": receive,
					})
					return
				}
				//  add conversation to pool
				conversation := &types.Conversation{
					Id:           response.ConversationId,
					ConnectionId: connection.Id,
				}
				conversationPool.Set(conversation)
				// Send response
				c.JSON(200, response)
				// Heartbeat
				connection.Heartbeat = time.Now()
				return
			} else {
				// error 500
				c.JSON(500, gin.H{
					"error": "Failed to find response from the client",
				})
				return
			}
		}
	}
}

func API_getConnections(c *gin.Context) {
	// Get connections
	var conections []*types.Connection
	connectionPool.Mu.RLock()
	for _, connection := range connectionPool.Connections {
		conections = append(conections, connection)
	}
	connectionPool.Mu.RUnlock()
	// Send connections
	c.JSON(200, gin.H{
		"connections": conections,
	})
}

func ping(connection_id string) bool {
	// Get connection
	connection, ok := connectionPool.Get(connection_id)
	// Send ping to the connection
	if ok {
		id := utils.GenerateId()
		send := types.Message{
			Id:      id,
			Message: "API-CGPT Ping",
		}
		connection.Ws.SetReadDeadline(time.Now().Add(5 * time.Second))
		err := connection.Ws.WriteJSON(send)
		if err != nil {
			return false
		}
		// wait for response with timeout
		for {
			// read message
			var receive types.Message
			err = connection.Ws.ReadJSON(&receive)
			if err != nil {
				return false
			}
			// check if message is the response
			if receive.Id == send.Id {
				return true
			} else {
				return false
			}
		}
	}
	return false
}
