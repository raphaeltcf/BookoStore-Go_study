package router

import (
	"github.com/gin-gonic/gin"
	"bookstore/internal/handlers"
)

func Setup() *gin.Engine {
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	v1 := r.Group("/v1")
	{
		b := v1.Group("/books")
		{
			b.POST("", handlers.CreateBook)
			b.GET("", handlers.ListBooks)
			b.GET("/:id", handlers.GetBook)
			b.PATCH("/:id", handlers.UpdateBook)
			b.DELETE("/:id", handlers.DeleteBook)
		}

		t := v1.Group("/rentals")
		{
			t.POST("", handlers.CreateRental)
			t.POST("/:id/return", handlers.ReturnRental)
			t.GET("", handlers.ListRentals)
			t.GET("/:id", handlers.GetRental)
		}
	}

	return r
}
