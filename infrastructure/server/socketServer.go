package server

import (
	"context"
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

	// Escuchar cambios en MongoDB
	wsServer.listenForSessionUpdates()
}

func (wsServer *WebSocketServer) listenForSessionUpdates() {
	ctx, cancel := context.WithCancel(context.Background()) // Usamos un contexto con cancelación
	defer cancel()

	collection := wsServer.mongoClient.Database("shufflerion").Collection("session")

	// Crear el change stream para detectar actualizaciones
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

		// Verificar si el campo "guest" fue actualizado y ya no está vacío
		if change.OperationType == "update" {
			if guest, ok := change.UpdateDescription.UpdatedFields["guest"]; ok {
				// Verificamos que "guest" es un objeto (no un array)
				if user, ok := guest.(map[string]interface{}); ok {
					// Aquí puedes acceder a las propiedades del usuario, por ejemplo, el email
					if email, ok := user["email"].(string); ok {
						log.Printf("Usuario en guest: %s", email)
					}
				} else {
					log.Printf("El campo 'guest' no es un objeto válido para sessionId %s", change.FullDocument.SessionID)
				}
				// Notificamos al cliente con la sesión actualizada
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
		"event": "guest__joined",
		"data":  session,
	}

	err := wsServer.conn.WriteJSON(message)
	if err != nil {
		log.Printf("Error enviando mensaje al cliente: %v", err)
	}
}
