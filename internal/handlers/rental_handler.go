package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"bookstore/internal/db"
	"bookstore/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateRental(c *gin.Context) {
	var in models.CreateRentalInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var b models.Book
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&b, in.BookID).Error; err != nil {
			return err
		}
		if b.AvailableCount <= 0 {
			return gin.Error{Type: gin.ErrorTypePublic, Meta: "livro indisponível"}
		}

		b.AvailableCount -= 1
		if err := tx.Save(&b).Error; err != nil {
			return err
		}

		r := models.Rental{
			BookID:     b.ID,
			RenterName: in.RenterName,
			Status:     models.RentalActive,
			RentedAt:   time.Now(),
			DueAt:      in.DueAt,
		}
		return tx.Create(&r).Error
	})

	if err != nil {
		if ge, ok := err.(gin.Error); ok && ge.Type == gin.ErrorTypePublic {
			c.JSON(http.StatusConflict, gin.H{"error": ge.Meta})
			return
		}
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "livro não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro criando aluguel"})
		return
	}

	c.Status(http.StatusCreated)
}

func ReturnRental(c *gin.Context) {
	id := c.Param("id")

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var r models.Rental
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("Book").First(&r, id).Error; err != nil {
			return err
		}

		if r.Status != models.RentalActive {
			return gin.Error{Type: gin.ErrorTypePublic, Meta: "apenas alugueis ACTIVE podem ser devolvidos"}
		}

		now := time.Now()
		check := now.Add(48 * time.Hour)

		r.Status = models.RentalReturnedPending
		r.ReturnedAt = &now
		r.CheckExpiresAt = &check

		return tx.Save(&r).Error
	})

	if err != nil {
		if ge, ok := err.(gin.Error); ok && ge.Type == gin.ErrorTypePublic {
			c.JSON(http.StatusConflict, gin.H{"error": ge.Meta})
			return
		}
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "aluguel não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro processando devolução"})
		return
	}

	c.Status(http.StatusNoContent)
}

func ListRentals(c *gin.Context) {
	var items []models.Rental
	if err := db.DB.Preload("Book").Order("id DESC").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro listando alugueis"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func GetRental(c *gin.Context) {
	var r models.Rental
	if err := db.DB.Preload("Book").First(&r, c.Param("id")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "aluguel não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro buscando aluguel"})
		return
	}
	c.JSON(http.StatusOK, r)
}
