package models

import "time"

type User struct {
	ID             uint
	Email          string
	Username       string
	PasswordHashed string
}

type Post struct {
	ID        uint
	Body      string
	CreatedAt time.Time
    Author    User      `gorm:"foreignKey:author"`
}
