package models

import "time"

type User struct {
    ID             uint     `json:"id"`
    Email          string   `json:"email"`
    Username       string   `json:"username"`
    PasswordHashed string   `json:"-"`
}

type Post struct {
    ID        uint      `json:"id"`
    Body      string    `json:"body"`
    CreatedAt time.Time `json:"created_at"`
    Author    User      `gorm:"foreignKey:author" json:"author"`
}
