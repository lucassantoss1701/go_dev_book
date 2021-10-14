package models

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID         uint64    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Nick       string    `json:"nick,omitempty"`
	Email      string    `json:"email,omitempty"`
	Passoword  string    `json:"password,omitempty"`
	Created_at time.Time `json:"Created_at,omitempty"`
}

func (user *User) Prepare(step string) error {
	if err := user.validate(step); err != nil {
		return err
	}
	user.format()
	return nil
}

func (user *User) validate(step string) error {
	if user.Name == "" {
		return errors.New("O nome é obrigátorio e não pode estar em branco")
	}
	if user.Nick == "" {
		return errors.New("O nick é obrigátorio e não pode estar em branco")
	}
	if user.Email == "" {
		return errors.New("O Email é obrigátorio e não pode estar em branco")
	}
	if step == "registration" && user.Passoword == "" {
		return errors.New("A senha é obrigátorio e não pode estar em branco")
	}
	return nil
}

func (user *User) format() {
	user.Name = strings.TrimSpace(user.Name)
	user.Nick = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Name)
}
