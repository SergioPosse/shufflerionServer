package domain

import "context"

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type User struct {
	Email  string
	Tokens Tokens
}

type Session struct {
	Id        string
	Host      User
	Guest     []User
	CreatedAt string
	UpdatedAt string
}

type CreateSession struct {
	Id string
	Host User
}

type UpdateSession struct {
	Id string
	Guest User
	UpdatedAt string
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session CreateSession) error
	GetSessionById(ctx context.Context, id string) (*Session, error)
	UpdateSession(ctx context.Context, session UpdateSession) (*Session, error)
}
