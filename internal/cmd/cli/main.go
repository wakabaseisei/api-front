package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/spf13/cobra"
)

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

func runMigration(db *sql.DB, dbName, direction string) error {
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

	if direction == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return err
		}
	} else if direction == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			return err
		}
	}
	return nil
}

var rootCmd = &cobra.Command{
	Use:   "migrate-cli",
	Short: "Database migration tool",
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		rdsHost := os.Getenv("RDS_HOST")
		dbName := os.Getenv("DB_NAME")
		iamUser := os.Getenv("DB_USER")
		region := os.Getenv("AWS_REGION")
		port := os.Getenv("DB_PORT")

		token, err := generateAuthToken(rdsHost, iamUser, region)
		if err != nil {
			log.Fatalf("generate IAM auth token: %v", err)
		}

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=true&multiStatements=true",
			iamUser, token, rdsHost, port, dbName)

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Failed to connect to DB: %v", err)
		}
		defer db.Close()

		if err := runMigration(db, dbName, "up"); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}

		fmt.Println("Migration completed successfully!")
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Revert the last migration",
	Run: func(cmd *cobra.Command, args []string) {
		rdsHost := os.Getenv("RDS_HOST")
		dbName := os.Getenv("DB_NAME")
		iamUser := os.Getenv("DB_USER")
		region := os.Getenv("AWS_REGION")
		port := os.Getenv("DB_PORT")

		token, err := generateAuthToken(rdsHost, iamUser, region)
		if err != nil {
			log.Fatalf("generate IAM auth token: %v", err)
		}

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=true&multiStatements=true",
			iamUser, token, rdsHost, port, dbName)

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Failed to connect to DB: %v", err)
		}
		defer db.Close()

		if err := runMigration(db, dbName, "down"); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}

		fmt.Println("Migration rollback completed!")
	},
}

func main() {
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
