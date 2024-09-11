package services

import (
    "fmt"
    "time"
)

type UserService struct {}

func NewUserService() *UserService {
    return &UserService{}
}

// GenerateUserAccessLink genera un enlace de acceso único para el usuario
func (s *UserService) GenerateUserAccessLink() string {
    // Aquí puedes implementar una lógica más robusta, como tokens JWT, etc.
    timestamp := time.Now().Unix()
    return fmt.Sprintf("https://myapp.com/access?user=ts=%d", timestamp)
}