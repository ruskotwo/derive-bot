package cmd

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"

	"github.com/ruskotwo/derive-bot/internal/config"
)

//go:embed migrations/*
var migrationsEmbed embed.FS

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Start Migration")

		mysqlDatabaseConfig := config.NewMysqlDatabaseConfig()

		db, err := sql.Open("mysql", mysqlDatabaseConfig.GetDsn())
		if err != nil {
			log.Fatalf("failed to connect to db: %v", err)
		}

		defer func(db *sql.DB) { _ = db.Close() }(db)

		db.SetMaxOpenConns(mysqlDatabaseConfig.MaxOpenConnections)
		db.SetMaxIdleConns(mysqlDatabaseConfig.MaxIdleConnections)
		db.SetConnMaxLifetime(mysqlDatabaseConfig.ConnectionMaxLifeTime)
		db.SetConnMaxIdleTime(mysqlDatabaseConfig.ConnectionMaxIdleTime)

		if len(args) == 0 {
			log.Fatalf("action is not provided")
		}

		action := args[0]
		options := args[1:]

		goose.SetBaseFS(migrationsEmbed)

		if err := goose.SetDialect(string(goose.DialectMySQL)); err != nil {
			log.Fatalf("failed to set dialect: %v", err)
		}

		err = goose.RunContext(cmd.Context(), action, db, "migrations", options...)
		if err != nil {
			log.Fatalf("failed to run migration: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
