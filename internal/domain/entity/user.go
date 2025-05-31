package entity

import (
	"time"
)

type User struct {
	ID        string
	Email     string
	Password  string
	Name      string
	Surname   string
	Phone     string
	Telegram  string
	IsAdmin   bool
	CreatedAt time.Time
}
