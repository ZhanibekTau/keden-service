package main

import (
	"keden-service/back/internal/app"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/sirupsen/logrus"
)

func init() {
	location, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		logrus.Fatalf("error loading location: %v", err.Error())
	}
	time.Local = location
}

func main() {
	newApp, err := app.NewApp()
	if err != nil {
		log.Fatalf("App error: %v", err)
	}

	var wg sync.WaitGroup

	// Start RabbitMQ consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := newApp.RunConsumer(); err != nil {
			log.Printf("Consumer error: %v", err)
		}
	}()

	if err = newApp.RunRestServer(); err != nil {
		log.Printf("REST server error: %v", err)
	}

	log.Println("Application started successfully")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gracefully...")
	wg.Wait()
	log.Println("All components stopped.")
}
