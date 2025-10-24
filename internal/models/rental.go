package models

import "time"

type RentalStatus string

const (
	RentalActive          RentalStatus = "ACTIVE"
	RentalReturnedPending RentalStatus = "RETURNED_PENDING_CHECK"
	RentalCompleted       RentalStatus = "COMPLETED"
)

type Rental struct {
	ID             uint         `json:"id" gorm:"primaryKey"`
	BookID         uint         `json:"book_id" binding:"required" gorm:"not null"`
	Book           Book         `json:"book" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	RenterName     string       `json:"renter_name" binding:"required,min=2,max=100" gorm:"not null"`
	Status         RentalStatus `json:"status" gorm:"index"`
	RentedAt       time.Time    `json:"rented_at"`
	DueAt          *time.Time   `json:"due_at,omitempty"`
	ReturnedAt     *time.Time   `json:"returned_at,omitempty"`
	CheckExpiresAt *time.Time   `json:"check_expires_at,omitempty"`
	CreatedAt      time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateRentalInput struct {
	BookID     uint       `json:"book_id" binding:"required,gt=0"`
	RenterName string     `json:"renter_name" binding:"required,min=2,max=100"`
	DueAt      *time.Time `json:"due_at"`
}
