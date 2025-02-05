package application

import (
	"context"
	"shufflerion/modules/session/domain"
)

type SessionUseCase struct {
	repo domain.SessionRepository
}

func NewSessionUseCase(repo domain.SessionRepository) *SessionUseCase {
	return &SessionUseCase{repo: repo}
}

func (uc *SessionUseCase) CreateSession(ctx context.Context, session domain.Session) error {
	return uc.repo.CreateSession(ctx, session)
}

func (uc *SessionUseCase) GetSession(ctx context.Context, id string) (*domain.Session, error) {
	return uc.repo.GetSessionById(ctx, id)
}

func (uc *SessionUseCase) UpdateSession(ctx context.Context, session domain.UpdateSession) (*domain.Session, error) {
	return uc.repo.UpdateSession(ctx, session)
}
