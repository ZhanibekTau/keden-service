package main

import (
	"keden-service/back/internal/app"
	"log"
	"os"
	"os/signal"
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

	if err = newApp.RunRestServer(); err != nil {
		log.Printf("REST server error: %v", err)
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gracefully...")
}
