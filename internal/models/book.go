package models

import "time"

type Book struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	Title          string    `json:"title"  binding:"required,min=2,max=100" gorm:"not null"`
	Author         string    `json:"author" binding:"required,min=2,max=100" gorm:"not null"`
	Price          float64   `json:"price"  binding:"required,gte=0"          gorm:"not null"`
	TotalCount     int       `json:"total_count"     gorm:"not null;default:0"`
	AvailableCount int       `json:"available_count" gorm:"not null;default:0"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateBookInput struct {
	Title      string  `json:"title"  binding:"required,min=2,max=100"`
	Author     string  `json:"author" binding:"required,min=2,max=100"`
	Price      float64 `json:"price"  binding:"required,gte=0"`
	TotalCount int     `json:"total_count" binding:"required,gte=0"`
}

type UpdateBookInput struct {
	Title      *string  `json:"title"       binding:"omitempty,min=2,max=100"`
	Author     *string  `json:"author"      binding:"omitempty,min=2,max=100"`
	Price      *float64 `json:"price"       binding:"omitempty,gte=0"`
	TotalCount *int     `json:"total_count" binding:"omitempty,gte=0"`
}
