package application

import (
    "shufflerion/infraestructure/services"
)

type GenerateUserAccessLinkUseCase struct {
    UserService *services.UserService
}

func NewGenerateUserAccessLinkUseCase(userService *services.UserService) *GenerateUserAccessLinkUseCase {
    return &GenerateUserAccessLinkUseCase{UserService: userService}
}

// Execute toma un userID y genera el enlace de acceso correspondiente
func (uc *GenerateUserAccessLinkUseCase) Execute() string {
    return uc.UserService.GenerateUserAccessLink()
}