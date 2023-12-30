package dto

type GetJwtInput struct {
	Email string `bson:"email" json:"email"`
	Senha string `bson:"-" json:"password,omitempty"`
}
type GetJWTOutput struct {
	AccessToken string `json:"access_token"`
}
