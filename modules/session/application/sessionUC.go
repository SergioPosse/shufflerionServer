package application

import (
	"context"
	"shufflerion/modules/session/domain"
)

type SessionUseCase struct {
	sessionRepository domain.SessionRepository
	spotifyService domain.SpotifyService
}

func NewSessionUseCase(repository domain.SessionRepository, spotify domain.SpotifyService) *SessionUseCase {
	return &SessionUseCase{sessionRepository: repository, spotifyService: spotify}
}

func (uc *SessionUseCase) CreateSession(ctx context.Context, session domain.Session) error {
	err := uc.sessionRepository.CreateSession(ctx, session)
	if err != nil {
		return err
	}

	res , err := uc.spotifyService.AddUserToApp(session.Host.Tokens.AccessToken, session.Host.Email)
	if err != nil || res {
		return err
	}
	return err
}

func (uc *SessionUseCase) GetSession(ctx context.Context, id string) (*domain.Session, error) {
	return uc.sessionRepository.GetSessionById(ctx, id)
}

func (uc *SessionUseCase) UpdateSession(ctx context.Context, session domain.UpdateSession) (*domain.Session, error) {
	return uc.sessionRepository.UpdateSession(ctx, session)
}
