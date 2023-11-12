package model

import (
	"log"
	"strings"
	"time"

	"github.com/katana/back-end/construcomprafacil-go/internal/security"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	DataType       string             `bson:"data_type" json:"-"`
	Username       string             `bson:"username" json:"username"`
	Name           string             `bson:"name" json:"name"`
	Password       string             `bson:"-" json:"password,omitempty"`
	HashedPassword string             `bson:"password" json:"-"`
	Email          string             `bson:"email" json:"email"`
	Enable         bool               `bson:"enable" json:"enable"`
	IsLocked       bool               `bson:"is_locked" json:"isLocked"`
	SuperAdmin     bool               `bson:"super_admin" json:"super_admin"`
	ChangePassword bool               `bson:"change_password" json:"change_password"`
	FirstAcccess   bool               `bson:"first_access" json:"first_access"`
	ExpireAt       string             `bson:"expire_at" json:"expire_at,omitempty"`
	CreatedAt      string             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt      string             `bson:"updated_at" json:"updated_at,omitempty"`
	DeletedAt      string             `bson:"deleted_at" json:"deleted_at,omitempty"`
}

type FilterUser struct {
	Username string
	Email    string
	Enable   string
}

func (u *User) Formatusr() error {
	u.Name = strings.TrimSpace(u.Name)
	u.Username = strings.TrimSpace(u.Username)
	u.Email = strings.TrimSpace(u.Email)

	if u.DataType == "create" {
		passwordHash, err := security.HashedPassword(u.Password)
		if err != nil {
			return err
		}
		u.HashedPassword = string(passwordHash)
		return nil
	}
	return nil
}

func (u *User) passwordToHash() {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
		if err != nil {
			log.Println("Erro to SetPassWord", err.Error())
		}

		u.HashedPassword = string(hashedPassword)
	}
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	if err != nil {
		log.Println("Erro to CheckPassword", err.Error())
		return false
	}
	return true
}

func (u *User) PrepareToSave() {
	dt := time.Now().Format(time.RFC3339)
	u.passwordToHash()
	u.DataType = "user"

	if u.ID.IsZero() {
		u.CreatedAt = dt
		u.UpdatedAt = dt
	} else {
		u.UpdatedAt = dt
	}
}
