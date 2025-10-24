package models

import "time"

type Book struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" binding:"required, min=2, max=100" gorm:"not null"`
	Author    string    `json:"author" binding:"required, min=2, max=100" gorm:"not null"`
	Price     float64   `json:"price" binding:"required, gte=0" gorm:"not null"`
	Available bool      `json:"available" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateBookRequest struct {
	Title     string  `json:"title" binding:"required, min=2, max=100"`
	Author    string  `json:"author" binding:"required, min=2, max=100"`
	Price     float64 `json:"price" binding:"required, gte=0"`
}

type UpdateBookRequest struct {
	Title     string  `json:"title" binding:"required, min=2, max=100"`
	Author    string  `json:"author" binding:"required, min=2, max=100"`
	Price     float64 `json:"price" binding:"required, gte=0"`
}