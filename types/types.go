package types

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Id      string `json:"id"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type ChatGptResponse struct {
	Id             string `json:"id"`
	ResponseId     string `json:"response_id"`
	ConversationId string `json:"conversation_id"`
	Content        string `json:"content"`
	Error          string `json:"error"`
}

type ChatGptRequest struct {
	MessageId      string `json:"message_id"`
	ConversationId string `json:"conversation_id"`
	ParentId       string `json:"parent_id"`
	Content        string `json:"content"`
}

type Connection struct {
	Ws              *websocket.Conn
	Id              string
	Heartbeat       time.Time
	LastMessageTime time.Time
}

type ConnectionPool struct {
	Connections map[string]*Connection
	Mu          sync.Mutex
}

// Get string with parameter id of a reference to connection using RW Mutex that returns
// a pointer to a connection
func (p *ConnectionPool) Get(id string) (*Connection, bool) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	conn, ok := p.Connections[id]
	// check if the connection  equal to nil
	if conn == nil {
		ok = false
	}
	return conn, ok
}

// set function with pointer to connection with ConnectionPool as receiver
func (p *ConnectionPool) Set(conn *Connection) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	p.Connections[conn.Id] = conn
}

// delete function with parameter id of a reference to connection using RW Mutex that returns error
func (p *ConnectionPool) Delete(id string) error {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	// delete connection from map
	delete(p.Connections, id)
	return nil
}

// NewConnectionPool function that returns a pointer to a new connection pool
func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{
		Connections: make(map[string]*Connection),
	}
}

//  Conversation struct with fields of Id and ConnectionId
type Conversation struct {
	Id           string
	ConnectionId string
}

// ConversationPool struct with fields of Conversations and Mu
type ConversationPool struct {
	Conversations map[string]*Conversation
	Mu            sync.Mutex
}

// Get function with parameter id of a reference to connection using RW Mutex that returns
// a pointer to a connection
func (p *ConversationPool) Get(id string) (*Conversation, bool) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	conversation, ok := p.Conversations[id]
	// check if the connection  equal to nil
	if conversation == nil {
		ok = false
	}
	return conversation, ok
}

// set function with pointer to connection with ConnectionPool as receiver
func (p *ConversationPool) Set(conversation *Conversation) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	p.Conversations[conversation.Id] = conversation
}

// delete function with parameter id of a reference to connection using RW Mutex that returns error
func (p *ConversationPool) Delete(id string) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	// delete connection from map
	delete(p.Conversations, id)
}

// NewConversationPool function that returns a pointer to a new connection pool
func NewConversationPool() *ConversationPool {
	return &ConversationPool{
		Conversations: make(map[string]*Conversation),
	}
}
