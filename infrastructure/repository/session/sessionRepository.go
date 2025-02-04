package repository

import (
	"context"
	"fmt"
	"shufflerion/modules/session/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoSessionRepository struct {
	collection *mongo.Collection
}

func NewMongoSessionRepository(db *mongo.Database) *MongoSessionRepository {
	return &MongoSessionRepository{
		collection: db.Collection("session"),
	}
}

// CreateSession: Crea una nueva sesión con el host y guest vacíos o con solo el host.
func (r *MongoSessionRepository) CreateSession(ctx context.Context, session domain.CreateSession) error {
	_, err := r.collection.InsertOne(ctx, session)
	return err
}

// GetSessionById: Obtiene una sesión por ID.
func (r *MongoSessionRepository) GetSessionById(ctx context.Context, id string) (*domain.Session, error) {
	var session domain.Session
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *MongoSessionRepository) UpdateSession(ctx context.Context, session domain.UpdateSession) (*domain.Session, error) {
	// Validar que session.Guest no sea nil
	if session.Guest.Email == "" {
		return nil, fmt.Errorf("Guest cannot be empty")
	}

	// Buscar la sesión existente por ID
	filter := bson.M{"id": session.Id}

	update := bson.M{
		"$push": bson.M{
			"guest": session.Guest,
		},
		"$setOnInsert": bson.M{ // Asegura que el campo existe si no está
			"guest": []domain.User{},
		},
	}

	// Realizar la actualización
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Verificar si se actualizó algo
	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("no session found with ID %v", session.Id)
	}

	// Obtener la sesión actualizada desde la base de datos
	var updatedSession domain.Session
	err = r.collection.FindOne(ctx, filter).Decode(&updatedSession)
	if err != nil {
		return nil, fmt.Errorf("error retrieving updated session: %v", err)
	}

	// Devolver la sesión actualizada
	return &updatedSession, nil
}
