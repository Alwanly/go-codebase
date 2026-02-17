package model

import "time"

// User model
type User struct {
	ID        string    `gorm:"primaryKey;column:id;type:varchar(32);not null" `
	Username  string    `gorm:"column:username;type:varchar(255);not null" `
	Password  string    `gorm:"column:password;type:varchar(255);not null" `
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null" `
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null" `
}

// TableName for User model
func (User) TableName() string {
	return "users"
}

// Users model
type Users []User
