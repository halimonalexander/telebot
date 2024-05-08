package storage

import (
	"crypto/sha1"
	"io"
	"telebot/lib/e"
	"time"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(username string) (*Page, error)
	Remove(*Page) error
	Prune(username string) error
	Exists(*Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
	Created  time.Time
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.WrapError("Unable to hash page", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.WrapError("Unable to hash page", err)
	}

	return string(h.Sum(nil)), nil
}
