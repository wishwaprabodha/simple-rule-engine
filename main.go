package main

import (
	"context"
	"fmt"
	"github.com/honestbank/tech-assignment-backend-engineer/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	env "github.com/joho/godotenv"

	"github.com/honestbank/tech-assignment-backend-engineer/controllers"
)

const envFile = ".env"

var loadEnv = env.Load

func run() (s *http.Server) {
	err := loadEnv(envFile)
	if err != nil {
		log.Fatal(err)
	}
	port, exist := os.LookupEnv("PORT")
	if !exist {
		log.Fatal("no port specified")
	}
	port = fmt.Sprintf(":%s", port)

	config.ConfigInstance.LoadConfig()
	mux := http.NewServeMux()

	mux.HandleFunc("/process", controllers.ProcessData)
	mux.HandleFunc("/update-config", config.UpdateConfigHandler)
	mux.HandleFunc("/get-config", config.GetConfigHandler)

	s = &http.Server{
		Addr:           port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        mux,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error listening on port: %s\n", err)
		}
	}()

	return
}

func main() {
	s := run()
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shut down")
	}
	log.Println("server exiting")
}
