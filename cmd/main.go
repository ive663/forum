package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ive663/forum/internal/delivery"
	"github.com/ive663/forum/internal/repository"
	"github.com/ive663/forum/internal/server"
	"github.com/ive663/forum/internal/service"
)

func main() {
	db, err := repository.Init()
	if err != nil {
		log.Print(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Can't close db err: %v\n", err)
		} else {
			log.Print("db closed")
		}
	}()
	if err := repository.CreateDatabase(db); err != nil {
		log.Print(err)
		return
	}
	repositories := repository.NewRepository(db)
	services := service.NewServices(repositories)
	handlers := delivery.NewHandler(services)
	server := new(server.Server)
	go func() {
		if err := server.Start(":8080", handlers.Handlers()); err != nil {
			log.Println(err)
			return
		}
	}()
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			if err := services.Auth.DeleteExpiredSessions(); err != nil {
				log.Println(err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		log.Print(err)
		return
	}
}
