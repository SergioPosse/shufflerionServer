package application

import (
	"shufflerion/modules/auth/domain"
	"fmt"
	"shufflerion/infrastructure/services"
)

type GetAccessTokensUseCase struct {
    AuthService *services.AuthService
}

func NewGetAccessTokensUseCase(authService *services.AuthService) *GetAccessTokensUseCase {
    return &GetAccessTokensUseCase{AuthService: authService}
}

func (uc *GetAccessTokensUseCase) Execute(code1 string, code2 string) []shared.GetAccessTokensResponse {
    data, err := uc.AuthService.GetAccessTokens(code1, code2)
    if(err != nil) {
        fmt.Print((err))
    }

    return 	data
}