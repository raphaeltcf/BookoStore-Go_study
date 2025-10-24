package jobs

import (
	"context"
	"log"
	"time"

	"bookstore/internal/db"
	"bookstore/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func StartCheckinJob(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	log.Println("[job] verificação de devoluções iniciado")
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("[job] verificação de devoluções encerrado")
				return
			case <-ticker.C:
				process()
			}
		}
	}()
}

func process() {
	now := time.Now()

	var rentals []models.Rental
	if err := db.DB.
		Where("status = ? AND check_expires_at IS NOT NULL AND check_expires_at <= ?", models.RentalReturnedPending, now).
		Find(&rentals).Error; err != nil {
		log.Printf("[job] erro buscando rentals: %v\n", err)
		return
	}

	for _, r := range rentals {
		err := db.DB.Transaction(func(tx *gorm.DB) error {
			var rr models.Rental
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Preload("Book").
				First(&rr, r.ID).Error; err != nil {
				return err
			}

			if rr.Status != models.RentalReturnedPending || rr.CheckExpiresAt == nil || rr.CheckExpiresAt.After(time.Now()) {
				return nil
			}

			rr.Book.AvailableCount += 1
			if err := tx.Save(&rr.Book).Error; err != nil {
				return err
			}

			rr.Status = models.RentalCompleted
			return tx.Save(&rr).Error
		})

		if err != nil {
			log.Printf("[job] erro processando rental %d: %v\n", r.ID, err)
		}
	}
}
