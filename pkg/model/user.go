package model

import (
	"encoding/json"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserInterface interface {
	String() string
}

type Usuario struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Nome      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Senha     string             `bson:"-" json:"password,omitempty"`
	Enable    bool               `bson:"enable" json:"enable"`
	CreatedAt string             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt string             `bson:"updated_at" json:"updated_at,omitempty"`
}

type FilterUsuario struct {
	Nome   string
	Email  string
	Enable string
}

func (u *Usuario) String() string {
	data, err := json.Marshal(u)

	if err != nil {
		log.Println("Error convert User to JSON")
		log.Println(err.Error())
		return ""
	}

	return string(data)
}

func (u *Usuario) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Senha), []byte(password))
	if err != nil {
		log.Println("Erro to CheckPassword", err.Error())
		return false
	}
	return true
}

func NewUsuario(nome, senha, email string) (*Usuario, error) {
	dt := time.Now().Format(time.RFC3339)
	if senha != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(senha), 10)
		if err != nil {
			log.Println("Erro to SetPassWord", err.Error())
			return nil, err
		}

		senha = string(hashedPassword)
	}
	tmp_user := &Usuario{
		Nome:      nome,
		Senha:     senha,
		Email:     email,
		Enable:    true,
		CreatedAt: dt,
		UpdatedAt: dt,
	}

	return tmp_user, nil
}
