package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wakabaseisei/api-front/internal/domain/repository"
	"github.com/wakabaseisei/api-front/internal/driver/grpc"
	infraRepo "github.com/wakabaseisei/api-front/internal/repository"
	"github.com/wakabaseisei/ms-protobuf/gen/go/ms/apifront/v1/apifrontv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// var (
// 	db     *sql.DB
// 	dbErr  error
// 	dbName = "DatabaseName"
// 	dbUser = "DatabaseUser"
// 	dbHost = "mysqlcluster.cluster-123456789012.us-east-1.rds.amazonaws.com"
// 	dbPort = 3306
// 	region = "us-east-1"
// )

func main() {
	services := repository.NewServices(infraRepo.NewUserRepository())
	service := grpc.NewAPIFrontService(services)
	mux := http.NewServeMux()

	path, handler := apifrontv1connect.NewGreetServiceHandler(service)

	mux.Handle(path, handler)
	mux.HandleFunc("/", healthCheckHandler)

	http.ListenAndServe(
		"0.0.0.0:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

// func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
// 	dbErr = checkAuroraConnection()
// 	if dbErr != nil {
// 		http.Error(w, fmt.Sprintf("Database connection failed: %v", dbErr), http.StatusServiceUnavailable)
// 		return
// 	}
// 	fmt.Fprintln(w, "OK")
// }

// func checkAuroraConnection() error {
// 	dbEndpoint := fmt.Sprintf("%s:%d", dbHost, dbPort)

// 	cfg, err := config.LoadDefaultConfig(context.TODO())
// 	if err != nil {
// 		return fmt.Errorf("configuration error: %w", err)
// 	}

// 	authToken, err := auth.BuildAuthToken(context.TODO(), dbEndpoint, region, dbUser, cfg.Credentials)
// 	if err != nil {
// 		return fmt.Errorf("create authentication token: %w", err)
// 	}

// 	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=true&allowCleartextPasswords=true",
// 		dbUser, authToken, dbEndpoint, dbName,
// 	)

// 	db, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		return fmt.Errorf("open database connection: %w", err)
// 	}
// 	defer db.Close()

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	if err := db.PingContext(ctx); err != nil {
// 		return fmt.Errorf("ping database: %w", err)
// 	}

// 	return nil
// }
