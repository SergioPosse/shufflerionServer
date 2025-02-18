package application_test

import (
	"context"
	"errors"
	"shufflerion/modules/session/application"
	"shufflerion/modules/session/domain"
	"testing"
  domainSong "shufflerion/modules/song/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSpotifyService struct {
	mock.Mock
}

func (m *MockSpotifyService) AddUserToApp(accessToken string, email string) (bool, error) {
	args := m.Called(accessToken, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockSpotifyService) FetchRandomSongs(accessToken string, quantity int) ([]domainSong.Song, error) {
	args := m.Called(accessToken, quantity)
	return args.Get(0).([]domainSong.Song), args.Error(1)
}



type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) CreateSession(ctx context.Context, session domain.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) GetSessionById(ctx context.Context, id string) (*domain.Session, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Session), args.Error(1)
}

func (m *MockSessionRepository) UpdateSession(ctx context.Context, session domain.UpdateSession) (*domain.Session, error) {
	args := m.Called(ctx, session)
	return args.Get(0).(*domain.Session), args.Error(1)
}

func TestCreateSession_Success(t *testing.T) {
  mockSpotify := new(MockSpotifyService)
  mockSpotify.On("AddUserToApp", mock.Anything, mock.Anything).Return(true, nil)
	mockRepo := new(MockSessionRepository)
	useCase := application.NewSessionUseCase(mockRepo, mockSpotify)

	session := domain.Session{Id: "123"}
	mockRepo.On("CreateSession", mock.Anything, session).Return(nil)

	err := useCase.CreateSession(context.TODO(), session)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateSession_Error(t *testing.T) {
  mockSpotify := new(MockSpotifyService)
  mockSpotify.On("AddUserToApp", mock.Anything, mock.Anything).Return(true, nil)

	mockRepo := new(MockSessionRepository)
	useCase := application.NewSessionUseCase(mockRepo, mockSpotify)

	session := domain.Session{Id: "123"}
	mockRepo.On("CreateSession", mock.Anything, session).Return(errors.New("DB error"))

	err := useCase.CreateSession(context.TODO(), session)

	assert.Error(t, err)
	assert.Equal(t, "DB error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetSession_Success(t *testing.T) {
  mockSpotify := new(MockSpotifyService)
  mockSpotify.On("AddUserToApp", mock.Anything, mock.Anything).Return(true, nil)

	mockRepo := new(MockSessionRepository)
	useCase := application.NewSessionUseCase(mockRepo, mockSpotify)

	expectedSession := &domain.Session{Id: "123"}
	mockRepo.On("GetSessionById", mock.Anything, "123").Return(expectedSession, nil)

	session, err := useCase.GetSession(context.TODO(), "123")

	assert.NoError(t, err)
	assert.Equal(t, expectedSession, session)
	mockRepo.AssertExpectations(t)
}

func TestGetSession_Error(t *testing.T) {
  mockSpotify := new(MockSpotifyService)
  mockSpotify.On("AddUserToApp", mock.Anything, mock.Anything).Return(true, nil)

	mockRepo := new(MockSessionRepository)
	useCase := application.NewSessionUseCase(mockRepo, mockSpotify)

	mockRepo.On("GetSessionById", mock.Anything, "123").Return((*domain.Session)(nil), errors.New("session not found"))

	session, err := useCase.GetSession(context.TODO(), "123")

	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Equal(t, "session not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdateSession_Success(t *testing.T) {
  mockSpotify := new(MockSpotifyService)
  mockSpotify.On("AddUserToApp", mock.Anything, mock.Anything).Return(true, nil)

	mockRepo := new(MockSessionRepository)
	useCase := application.NewSessionUseCase(mockRepo, mockSpotify)

	updateSession := domain.UpdateSession{Id: "123"}
	expectedSession := &domain.Session{Id: "123"}
	mockRepo.On("UpdateSession", mock.Anything, updateSession).Return(expectedSession, nil)

	session, err := useCase.UpdateSession(context.TODO(), updateSession)

	assert.NoError(t, err)
	assert.Equal(t, expectedSession, session)
	mockRepo.AssertExpectations(t)
}

func TestUpdateSession_Error(t *testing.T) {
  mockSpotify := new(MockSpotifyService)
  mockSpotify.On("AddUserToApp", mock.Anything, mock.Anything).Return(true, nil)

	mockRepo := new(MockSessionRepository)
	useCase := application.NewSessionUseCase(mockRepo, mockSpotify)

	updateSession := domain.UpdateSession{Id: "123"}
	mockRepo.On("UpdateSession", mock.Anything, updateSession).Return((*domain.Session)(nil), errors.New("update failed"))

	session, err := useCase.UpdateSession(context.TODO(), updateSession)

	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Equal(t, "update failed", err.Error())
	mockRepo.AssertExpectations(t)
}
