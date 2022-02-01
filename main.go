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
	"github.com/joho/godotenv"
	"github.com/manpreet1130/library-management/database"
	"github.com/manpreet1130/library-management/models"
	"github.com/manpreet1130/library-management/routes"
)

// init function used to initialize and connect to the database,
// create required tables in the database and add required foreign keys
func init() {
	database.Connect()
	db := database.GetDB()
	db.AutoMigrate(&models.User{}, &models.Book{}, &models.Cart{})
	db.Model(&models.Book{}).AddForeignKey("cart_uuid", "carts(id)", "RESTRICT", "CASCADE")
}

func main() {
	fmt.Println("Starting server..")

	if err := godotenv.Load(); err != nil {
		log.Fatal("Could not load env file")
	}

	router := mux.NewRouter()

	routes.UserRoutes(router)
	routes.BookRoutes(router)
	routes.CartRoutes(router)
	http.Handle("/", router)

	server := &http.Server{
		Addr:    os.Getenv("PORT"),
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
