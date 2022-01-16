package main

import (
	"context"
	"go-microservices/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Create a now logger
	logger := log.New(os.Stdout, "product-api", log.LstdFlags)

	// Create the handlers
	productsHandler := handlers.NewProducts(logger)

	// Create a new serve mux and register the handlers
	serveMux := http.NewServeMux()
	serveMux.Handle("/", productsHandler)

	// Create a new server
	server := &http.Server{
		Addr:         ":9090",           // Binding port / address:port
		Handler:      serveMux,          // Set the default handler
		ErrorLog:     logger,            // Set the logger for the server
		IdleTimeout:  120 * time.Second, // Max time for connection using TCP keep-alive
		ReadTimeout:  1 * time.Second,   // Max time to read request from the client
		WriteTimeout: 1 * time.Second,   // Max time to write response from the client
	}

	// Start the server on a go routine
	go func() {
		logger.Println("Starting server on port 9090!")

		// Start the server
		err := server.ListenAndServe()

		// Log the error
		if err != nil {
			logger.Fatal(err)
		}
	}()

	// Trap sigterm or interrupt and gracefully shutdown the server
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	sig := <-signalChannel
	logger.Println("Received terminate, graceful shutdown", sig)

	timeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second)

	server.Shutdown(timeoutContext)
}
