package controllers

import (
	"backend-in-go/db"
	"backend-in-go/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// In production, restrict origins.
	CheckOrigin: func(r *http.Request) bool { return true },
}
var message models.Message

type upcomingMessage struct {
	Type   string            
	Sender string
	Text   string 
	Reciever string
}

type Client struct {
	ID   string 
	Conn *websocket.Conn
	Send chan []byte
}
type Manager struct{
	mu      sync.RWMutex
	ClientList map[string]*Client
}
func NewManager() *Manager{
	return &Manager{
		ClientList: make(map[string]*Client),
	}
}

func (m *Manager) AddClient(c *Client){
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ClientList[c.ID] = c
}
func (m *Manager) RemoveClient(c *Client){
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.ClientList, c.ID)
}
func (m *Manager) GetClient(id string) (*Client, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.ClientList[id]
	return c, ok
}


func ServeWS(m *Manager,w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil{
		// w.WriteHeader(http.StatusInternalServerError)
		log.Println("Failed to upgrade connection", err)
		return
	}
	// _, data, err := conn.ReadMessage()
	// log.Println(string(data))
	// if err != nil{
	// 	log.Fatal("Failed to read message", err)
	// 	return
	// } 
	// Create a new client
    
	var RegisterMsg map[string]string
	// Unmarshal the JSON data into the RegisterMsg struct
	err = conn.ReadJSON(&RegisterMsg)
	if err != nil {
		log.Fatal("Failed to read message", err)
		return
	}
	fmt.Println(RegisterMsg)

	if RegisterMsg["Type"] != "register" || RegisterMsg["Sender"] == "" {
		log.Fatal("Failed to register", err)
		return
	}
	senderId, err := primitive.ObjectIDFromHex(RegisterMsg["Sender"])
    if err != nil {
        log.Fatal("Failed to convert sender ID to ObjectID", err)
        return
    }
	// recieverId,err := primitive.ObjectIDFromHex(RegisterMsg["Reciever"].(string))
	// if err != nil {
    //     log.Fatal("Failed to convert reciever ID to ObjectID", err)
    //     return
    // }

	message = models.Message{
		Type : RegisterMsg["Type"],
		Sender : senderId,
		// Reciever: recieverId,
		Text: RegisterMsg["Text"],
	}
	_, err = db.Collection_messages.InsertOne(context.TODO(), message)
    if err != nil {
		log.Fatal("Failed to insert message into Database", err)
	}
	// Create a new client
	client := &Client{
		ID:   RegisterMsg["sender"],
		Conn: conn,
		Send: make(chan []byte),
	}

	m.AddClient(client)

	go func(c *Client){
		defer func(){
           c.Conn.Close()
		   m.RemoveClient(c)
		}()

		for msg := range c.Send{
			err := c.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil{
				log.Println("Failed to send message", err)
				return
			}
		}
	}(client)
	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		var msg upcomingMessage
		if err := json.Unmarshal(payload, &msg); err != nil {
			client.Send <- []byte(`{"type":"error","text":"bad json"}`)
			continue
		}

		switch msg.Type {
		case "message":
			senderId, err := primitive.ObjectIDFromHex(msg.Sender)
    		if err != nil {
        	log.Fatal("Failed to convert sender ID to ObjectID", err)
        	return
    		}
			recieverId,err := primitive.ObjectIDFromHex(msg.Reciever)
			if err != nil {
      	    	log.Fatal("Failed to convert reciever ID to ObjectID", err)
        		return
    		}

			message = models.Message{
				Type: msg.Type,
				Sender : senderId,
				Reciever: recieverId,
				Text: msg.Text,
			}

			_, err = db.Collection_messages.InsertOne(context.TODO(), message)
			if err != nil {
				log.Fatal("Failed to insert message into Database", err)
			}



			if msg.Reciever == "" {
				client.Send <- []byte(`{"type":"error","text":"missing recipient"}`)
				continue
			}
			if msg.Sender == "" {
				// default to registered ID if not explicitly set
				msg.Sender = client.ID
			}
			target, ok := m.GetClient(msg.Reciever)
			if !ok {
				client.Send <- []byte(`{"type":"error","text":"recipient offline"}`)
					continue
			}
			forward, _ := json.Marshal(msg)
			select {
			case target.Send <- forward:
			default:
				client.Send <- []byte(`{"type":"error","text":"recipient busy"}`)
			}
		default:
			client.Send <- []byte(`{"type":"error","text":"unsupported type"}`)
		}
	}

	// Cleanup handled by writer goroutine defer
	close(client.Send)
}


