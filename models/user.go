package models

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"os"

	"github.com/guregu/null"
)

type Password string

type User struct {
	ID           int         `gorm:"column:id;primary_key" json:"id"`
	Presentation string      `gorm:"NOT NULL;column:presentation;" json:"presentation"`
	Login        string      `gorm:"NOT NULL;column:login;size:50" json:"login"`
	Password     Password    `gorm:"NOT NULL;column:password;size:32" json:"password,omitempty"`
	Email        null.String `gorm:"column:email;size:80" json:"email"`
	Phone        null.String `gorm:"column:phone;size:20" json:"phone"`
	Status       int         `gorm:"NOT NULL;column:status;type:tinyint(1);default:0" json:"-"`
	Role         int         `gorm:"NOT NULL;column:role;type:tinyint(1);default:0" json:"role"`
	Avatar       int
	Theme        int         `gorm:"NOT NULL;column:theme;type:tinyint(1);default:0" json:"theme"`
	ThemeColor   null.String `gorm:"column:theme_color;size:80" json:"theme_color"`
}

type Users []*User

func GetEncryptPassword(pass string) string {
	salt := os.Getenv("password_salt")
	hash := md5.New()
	hash.Write([]byte(pass + salt))
	return hex.EncodeToString(hash.Sum(nil))
}

func (p *Password) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*p = Password(s)
	return nil
}

func (p Password) String() string {
	return string(p)
}

func (p *Password) Encrypt() {
	pass := p.String()
	pass = GetEncryptPassword(pass)
	*p = Password(pass)
}
