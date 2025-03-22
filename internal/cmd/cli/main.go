package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
)

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	rdsHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	iamUser := os.Getenv("DB_USER")
	region := os.Getenv("AWS_REGION")
	port := os.Getenv("DB_PORT")
	dbEndpoint := fmt.Sprintf("%s:%s", rdsHost, port)

	token, err := generateAuthToken(dbEndpoint, iamUser, region)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("generate IAM auth token: %v", err),
		}, nil
	}

	caPath := "/etc/ssl/certs/rds-ca.pem"
	escapedCAPath := url.QueryEscape(caPath)

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=rds&x-tls-ca=%s&multiStatements=true",
		iamUser, token, dbEndpoint, dbName, escapedCAPath)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Failed to connect to DB: %v", err),
		}, nil
	}
	defer db.Close()

	if err := runMigration(db, dbName); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Migration failed: %v", err),
		}, nil
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "\"Migration was successful!\"",
	}
	return response, nil
}

func generateAuthToken(host, user, region string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return "", err
	}

	authenticationToken, terr := auth.BuildAuthToken(
		context.TODO(), host, region, user, cfg.Credentials)
	if terr != nil {
		return "", terr
	}

	return authenticationToken, nil
}

func runMigration(db *sql.DB, dbName string) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get executable path: %w", err)
	}
	migrationDir := filepath.Join(filepath.Dir(exePath), "../db/migrations")

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationDir),
		dbName,
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
