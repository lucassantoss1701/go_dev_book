package models

import (
	"errors"
	"strings"
	"time"
)

type Publish struct {
	ID         uint64    `json:"id,omitempty"`
	Title      string    `json:"title,omitempty"`
	Content    string    `json:"content,omitempty"`
	AuthorID   uint64    `json:"authorId,omitempty"`
	AuthorNick uint64    `json:"authorNick,omitempty"`
	Likes      uint64    `json:"likes"`
	Created_at time.Time `json:"created_at,omitempty"`
}

func (publish *Publish) Prepare() error {
	if erro := publish.validate(); erro != nil {
		return erro
	}

	publish.format()
	return nil
}

func (publish *Publish) validate() error {
	if publish.Title == "" {
		return errors.New("O título é obrigatório e não pode estar em branco")
	}

	if publish.Content == "" {
		return errors.New("O conteúdo é obrigatório e não pode estar em branco")
	}

	return nil
}

func (publish *Publish) format() {
	publish.Title = strings.TrimSpace(publish.Title)
	publish.Content = strings.TrimSpace(publish.Content)
}
