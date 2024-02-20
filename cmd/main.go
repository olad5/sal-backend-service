package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/olad5/sal-backend-service/internal/app/router"
)

func main() {
	port := os.Getenv("PORT")
	ctx := context.Background()
	appRouter := router.NewHttpRouter(ctx)
	server := &http.Server{Addr: ":" + port, Handler: appRouter}
	go func() {
		log.Printf("starting application server on  http://localhost:" + port + "\n")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("HTTP server ListenAndServe: %v", err)
		}
	}()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exiting gracefully")
}
