package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"bookstore/internal/db"
	"bookstore/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateBook(c *gin.Context) {
	var in models.CreateBookInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	b := models.Book{
		Title:          in.Title,
		Author:         in.Author,
		Price:          in.Price,
		TotalCount:     in.TotalCount,
		AvailableCount: in.TotalCount, 
	}

	if err := db.DB.Create(&b).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro criando livro"})
		return
	}
	c.JSON(http.StatusCreated, b)
}


func ListBooks(c *gin.Context) {
	var items []models.Book
	if err := db.DB.Order("id DESC").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro listando livros"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func GetBook(c *gin.Context) {
	var b models.Book
	if err := db.DB.First(&b, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "livro não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro buscando livro"})
		return
	}
	c.JSON(http.StatusOK, b)
}

func UpdateBook(c *gin.Context) {
	var in models.UpdateBookInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var b models.Book
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&b, c.Param("id")).Error; err != nil {
			return err
		}

		if in.Title != nil {
			b.Title = *in.Title
		}
		if in.Author != nil {
			b.Author = *in.Author
		}
		if in.Price != nil {
			b.Price = *in.Price
		}
		if in.TotalCount != nil {
			rentedNow := b.TotalCount - b.AvailableCount
			if *in.TotalCount < rentedNow {
				return gin.Error{Type: gin.ErrorTypePublic, Meta: "total_count não pode ser menor que itens alugados atualmente"}
			}
			b.TotalCount = *in.TotalCount
			b.AvailableCount = b.TotalCount - rentedNow
		}

		return tx.Save(&b).Error
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "livro não encontrado"})
			return
		}
		if ge, ok := err.(gin.Error); ok && ge.Type == gin.ErrorTypePublic {
			c.JSON(http.StatusConflict, gin.H{"error": ge.Meta})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro atualizando livro"})
		return
	}

	c.Status(http.StatusNoContent)
}

func DeleteBook(c *gin.Context) {
	var count int64
	db.DB.Model(&models.Rental{}).
		Where("book_id = ? AND status IN ?", c.Param("id"),
			[]models.RentalStatus{models.RentalActive, models.RentalReturnedPending}).
		Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "não é possível deletar livro com alugueis ativos/pending"})
		return
	}

	if err := db.DB.Delete(&models.Book{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro deletando livro"})
		return
	}
	c.Status(http.StatusNoContent)
}
