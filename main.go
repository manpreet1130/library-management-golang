package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/manpreet1130/library-management/database"
	"github.com/manpreet1130/library-management/models"
	"github.com/manpreet1130/library-management/routes"
)

func init() {
	database.Connect()
	db := database.GetDB()
	db.AutoMigrate(&models.User{}, &models.Book{})
}

func main() {
	fmt.Println("Starting server..")

	// port := os.Getenv("PORT")
	router := mux.NewRouter()

	routes.UserRoutes(router)
	routes.BookRoutes(router)
	http.Handle("/", router)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("could not instantiate server")
		}
	}()

	<-sigChan
	log.Println("Interrupt caught, commencing graceful shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Graceful shutdown failed")
	}
	log.Println("Shutdown successful")
}
