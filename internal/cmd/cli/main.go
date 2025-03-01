package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/spf13/cobra"
)

// IAM 認証用のトークン取得
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

// マイグレーションの実行
func runMigration(db *sql.DB, direction string) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // ローカルの `migrations` フォルダを使用
		"mysql",
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

// Cobra コマンドの定義
var rootCmd = &cobra.Command{
	Use:   "migrate-cli",
	Short: "Database migration tool",
}

// `migrate up` コマンド
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		// 環境変数から設定取得
		rdsHost := os.Getenv("RDS_HOST")
		dbName := os.Getenv("DB_NAME")
		iamUser := os.Getenv("DB_USER")
		region := os.Getenv("AWS_REGION")
		port := os.Getenv("DB_PORT")

		// IAM 認証のトークン取得
		token, err := generateAuthToken(rdsHost, iamUser, region)
		if err != nil {
			log.Fatalf("Failed to generate IAM auth token: %v", err)
		}

		// DB接続
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=true&allowCleartextPasswords=true",
			iamUser, token, rdsHost, port, dbName)

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Failed to connect to DB: %v", err)
		}
		defer db.Close()

		// マイグレーション実行
		if err := runMigration(db, "up"); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}

		fmt.Println("Migration completed successfully!")
	},
}

// `migrate down` コマンド
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Revert the last migration",
	Run: func(cmd *cobra.Command, args []string) {
		// 環境変数から設定取得
		rdsHost := os.Getenv("RDS_HOST")
		dbName := os.Getenv("DB_NAME")
		iamUser := os.Getenv("DB_USER")
		region := os.Getenv("AWS_REGION")

		// IAM 認証のトークン取得
		token, err := generateAuthToken(rdsHost, iamUser, region)
		if err != nil {
			log.Fatalf("Failed to generate IAM auth token: %v", err)
		}

		// DB接続
		dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?tls=true&multiStatements=true",
			iamUser, token, rdsHost, dbName)

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Failed to connect to DB: %v", err)
		}
		defer db.Close()

		// マイグレーション実行
		if err := runMigration(db, "down"); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}

		fmt.Println("Migration rollback completed!")
	},
}

func main() {
	// コマンドを登録
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)

	// 実行
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
