package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bookstore/internal/db"
	"bookstore/internal/jobs"
	"bookstore/internal/models"
	"bookstore/internal/router"
)

func main() {
	db.Connect("file:bookstore.db?cache=shared&_fk=1", &models.Book{}, &models.Rental{})


	r := router.Setup()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	jobCtx, cancelJobs := context.WithCancel(context.Background())
	jobs.StartCheckinJob(jobCtx)

	go func() {
		log.Println("servindo em http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("erro no servidor: %v", err)
		}
	}()


	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("encerrando...")
	cancelJobs()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown forçado: %v", err)
	}
	log.Println("bye 👋")
}
