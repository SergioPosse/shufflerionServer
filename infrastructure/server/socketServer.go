package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"shufflerion/modules/session/domain"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketServer struct {
	conn        *websocket.Conn
	mongoClient *mongo.Client
	sessionID   string
}

func NewWebSocketServer(mongoClient *mongo.Client) *WebSocketServer {
	return &WebSocketServer{mongoClient: mongoClient}
}

func (wsServer *WebSocketServer) HandleConnection(w http.ResponseWriter, r *http.Request) {
	var err error
	wsServer.conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error al actualizar a WebSocket: %v", err)
		return
	}
	defer wsServer.conn.Close()

	log.Println("Cliente conectado")

	// Leer el mensaje inicial para obtener el sessionId
	_, msg, err := wsServer.conn.ReadMessage()
	if err != nil {
		log.Printf("Error leyendo mensaje inicial: %v", err)
		return
	}

	var message struct {
		Action    string `json:"action"`
		SessionID string `json:"sessionId"`
	}
	if err := json.Unmarshal(msg, &message); err != nil {
		log.Printf("Error decodificando mensaje: %v", err)
		return
	}

	if message.Action == "subscribe" && message.SessionID != "" {
		wsServer.sessionID = message.SessionID
		log.Printf("Suscrito a actualizaciones de la sesión: %s", wsServer.sessionID)
		wsServer.listenForSessionUpdates()
	} else {
		log.Println("Mensaje de suscripción inválido")
	}
}

func (wsServer *WebSocketServer) listenForSessionUpdates() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	collection := wsServer.mongoClient.Database("shufflerion").Collection("session")

	changeStream, err := collection.Watch(ctx, mongo.Pipeline{}, options.ChangeStream().SetFullDocument(options.UpdateLookup))
	if err != nil {
		log.Fatalf("Error creando Change Stream: %v", err)
	}
	defer changeStream.Close(ctx)

	for changeStream.Next(ctx) {
		var change struct {
			OperationType string `bson:"operationType"`
			FullDocument  struct {
				SessionID string      `bson:"id"`
				Guest     domain.User `bson:"guest"`
			} `bson:"fullDocument"`
			UpdateDescription struct {
				UpdatedFields map[string]interface{} `bson:"updatedFields"`
			} `bson:"updateDescription"`
		}

		if err := changeStream.Decode(&change); err != nil {
			log.Printf("Error decodificando cambio: %v", err)
			continue
		}

		// Filtrar por el sessionId suscrito
		if change.FullDocument.SessionID == wsServer.sessionID && change.OperationType == "update" {
			if guest, ok := change.UpdateDescription.UpdatedFields["guest"]; ok {
				if user, ok := guest.(map[string]interface{}); ok {
					if email, ok := user["email"].(string); ok {
						log.Printf("Usuario en guest actualizado: %s", email)
					}
				}
				wsServer.notifyClient(change.FullDocument)
			}
		}
	}

	if err := changeStream.Err(); err != nil {
		log.Printf("Error en el Change Stream: %v", err)
	}
}

func (wsServer *WebSocketServer) notifyClient(session interface{}) {
	if wsServer.conn == nil {
		return
	}

	message := map[string]interface{}{
		"event": "guest_joined",
		"data":  session,
	}

	err := wsServer.conn.WriteJSON(message)
	if err != nil {
		log.Printf("Error enviando mensaje al cliente: %v", err)
	}
}