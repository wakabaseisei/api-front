package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"buf.build/gen/go/wakabaseisei/ms-protobuf/connectrpc/go/ms/apifront/v1/apifrontv1connect"
	connctcors "connectrpc.com/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/cors"
	"github.com/wakabaseisei/api-front/internal/config"
	"github.com/wakabaseisei/api-front/internal/domain/service"
	"github.com/wakabaseisei/api-front/internal/driver/client"
	"github.com/wakabaseisei/api-front/internal/driver/grpc"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	cfg, cerr := config.NewConfig()
	if cerr != nil {
		log.Fatalf("New Config: %v", cerr)
	}

	sgCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	services := service.NewServices(client.NewUserService(cfg.UserServiceEndpoint))
	service := grpc.NewAPIFrontService(services)
	mux := http.NewServeMux()

	path, handler := apifrontv1connect.NewGreetServiceHandler(service)

	mux.Handle(path, handler)
	mux.HandleFunc("/", healthCheckHandler)

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: withCORS(h2c.NewHandler(mux, &http2.Server{})),
	}

	done := make(chan error, 1)
	go func() {
		done <- server.ListenAndServe()
	}()

	select {
	case err := <-done:
		if err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	case <-sgCtx.Done():
		log.Println("Server stopping")
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(c); err != nil {
			log.Fatalf("HTTP server Shutdown: %v", err)
		}
		log.Println("Server gracefully stopped")
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

func withCORS(h http.Handler) http.Handler {
	middleware := cors.New(cors.Options{
		// TODO:
		AllowedOrigins: []string{"*"},
		AllowedMethods: connctcors.AllowedMethods(),
		AllowedHeaders: connctcors.AllowedHeaders(),
		ExposedHeaders: connctcors.ExposedHeaders(),
	})
	return middleware.Handler(h)
}
