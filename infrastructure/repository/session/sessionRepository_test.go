package repository_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	repository "shufflerion/infrastructure/repository/session"
	"shufflerion/modules/session/domain"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoSessionRepository(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		fmt.Println("Error loading .env")
	}

	mongoURI := fmt.Sprintf(
		"mongodb+srv://%s:%s@clusterfree.x3n59lo.mongodb.net/?retryWrites=true&w=majority&appName=ClusterFree",
		os.Getenv("DB_MONGO_USER"),
		os.Getenv("DB_MONGO_PASSWORD"),
	)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	assert.NoError(t, err)

	db := client.Database("testdb")
	repo := repository.NewMongoSessionRepository(db)

	db.Collection("session").Drop(context.TODO())
	t.Cleanup(func() {
		db.Collection("session").Drop(context.TODO())
	})

	t.Run("CreateSession should insert a session", func(t *testing.T) {
		session := domain.Session{Id: "123", Host: domain.User{Email: "host@test.com"}}
		session.Guest = domain.User{}
		err := repo.CreateSession(context.TODO(), session)
		assert.NoError(t, err)

		var foundSession domain.Session
		err = db.Collection("session").FindOne(context.TODO(), bson.M{"id": "123"}).Decode(&foundSession)
		assert.NoError(t, err)
		assert.Equal(t, "123", foundSession.Id)
		assert.Equal(t, "host@test.com", foundSession.Host.Email)
	})

	t.Run("GetSessionById should return the correct session", func(t *testing.T) {
		session := domain.Session{Id: "456", Host: domain.User{Email: "host2@test.com"}}
		session.Guest = domain.User{}
		_ = repo.CreateSession(context.TODO(), session)

		result, err := repo.GetSessionById(context.TODO(), "456")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "host2@test.com", result.Host.Email)
	})

	t.Run("UpdateSession should update guest list", func(t *testing.T) {
		session := domain.Session{Id: "789", Host: domain.User{Email: "host3@test.com"}}
		session.Guest = domain.User{}
		_ = repo.CreateSession(context.TODO(), session)

		update := domain.UpdateSession{Id: "789", Guest: domain.User{Email: "guest@test.com"}}
		updatedSession, err := repo.UpdateSession(context.TODO(), update)

		assert.NoError(t, err)
		assert.NotNil(t, updatedSession, "updatedSession es nil")

		if updatedSession != nil{
			assert.Equal(t, "guest@test.com", updatedSession.Guest.Email)
		} else {
			t.Errorf("the guest list has not be updated correctly")
		}
	})

	t.Run("UpdateSession should return error if guest is empty", func(t *testing.T) {
		update := domain.UpdateSession{Id: "999", Guest: domain.User{}}
		_, err := repo.UpdateSession(context.TODO(), update)
		assert.Error(t, err)
		assert.Equal(t, "guest cannot be empty", err.Error())
	})

	t.Run("UpdateSession should return error if session not found", func(t *testing.T) {
		update := domain.UpdateSession{Id: "000", Guest: domain.User{Email: "guest@test.com"}}
		_, err := repo.UpdateSession(context.TODO(), update)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no session found with ID")
	})
}
