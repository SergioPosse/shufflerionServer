package repository

import (
	"context"
	"fmt"
	"shufflerion/modules/session/domain"
	"time"

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

func (r *MongoSessionRepository) CreateSession(ctx context.Context, session domain.Session) error {
	if session.Guest == nil {
		session.Guest = []domain.User{}
		session.UpdatedAt = ""
	}
	session.CreatedAt = time.Now().Format(time.RFC3339)
	_, err := r.collection.InsertOne(ctx, session)
	return err
}

func (r *MongoSessionRepository) GetSessionById(ctx context.Context, id string) (*domain.Session, error) {
	var session domain.Session
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *MongoSessionRepository) UpdateSession(ctx context.Context, session domain.UpdateSession) (*domain.Session, error) {
	if session.Guest.Email == "" {
		return nil, fmt.Errorf("Guest cannot be empty")
	}

	filter := bson.M{"id": session.Id}

	update := bson.M{
		"$addToSet": bson.M{
			"guest": session.Guest,
		},
		"$set": bson.M{
			"updatedat": time.Now().Format(time.RFC3339),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("no session found with ID %v", session.Id)
	}

	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("guest with email %s already exists in the session", session.Guest.Email)
	}

	var updatedSession domain.Session
	err = r.collection.FindOne(ctx, filter).Decode(&updatedSession)
	if err != nil {
		return nil, fmt.Errorf("error retrieving updated session: %v", err)
	}

	return &updatedSession, nil
}

