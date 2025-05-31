package entity

import (
	"time"
)

type CardStatus string

const (
	StatusLost  CardStatus = "lost"
	StatusFound CardStatus = "found"
)

type Owner struct {
	ID       string
	Name     string
	Surname  string
	Phone    string
	Telegram string
}

type Card struct {
	ID          string
	Title       string
	Description string
	Latitude    float64
	Longitude   float64
	City        string
	Street      string
	PreviewURL  string
	Images      []string
	Status      CardStatus
	OwnerID     string
	CreatedAt   time.Time

	Owner     Owner
	DistanceM float64
}
