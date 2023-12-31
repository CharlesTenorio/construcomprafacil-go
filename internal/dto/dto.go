package dto

type GetJwtInput struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}
type GetJWTOutput struct {
	AccessToken string `json:"access_token"`
}
