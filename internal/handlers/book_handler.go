package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"bookstore/internal/models"
	"bookstore/internal/db"
	"gorm.io/gorm"
)

func CreateBook(c *gin.Context) {
	var in models.CreateBookRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := models.Book{
		Title: in.Title,
		Author: in.Author,
		Price: in.Price,
		Available: in.Price,
	}

	if err := db.DB.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"book": book})

}

func ListBooks(c *gin.Context) {
	var books []models.Book
	if err := db.DB.Order("id DESC").Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books": books})
}

func GetBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	if err := db.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"book": book})
}

func UpdateBook(c *gin.Context) {
	var in models.UpdateBookRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var book models.Book
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&book, c.Param("id")).Error; err != nil {
			return err 
		}

		if in.Title != nil {
			book.Title = *in.Title
		}
		if in.Author != nil {
			book.Author = *in.Author
		}
		if in.Price != nil {
			rentedNow := book.Price - book.Available 
			if *in.Price < rentedNow {
				return gin.Error {
					Err: errors.New("price cannot be less than the rented price"),
					Type: gin.ErrorTypePublic,	
					Meta: "total_count não pode ser menor que itens alugados atualmente",
				}
		}}
		book.available = *in.Price - rentedNow
		book.price = *in.Price
		if err := tx.Save(&book).Error; err != nil {
			return err
		}
	})

	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
	// regra simples: impedir deletar se há rentals ativos/pending
	var count int64
	db.DB.Model(&models.Rental{}).
		Where("book_id = ? AND status IN ?", c.Param("id"), []models.RentalStatus{models.RentalActive, models.RentalReturnedPending}).
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