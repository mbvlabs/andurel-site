package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"andurel-site/config"
	"andurel-site/migrations"

	"github.com/mbvlabs/andurel/pkg/storage"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("unexpected argument %q", args[0])
	}

	if err := config.LoadEnvironment(); err != nil {
		return err
	}

	ctx := context.Background()

	cfg, err := config.NewDatabase()
	if err != nil {
		return err
	}

	db, err := storage.NewPostgres(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	if err := storage.RunMigrations(ctx, db, migrations.Migrations, "."); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	fmt.Println("Migrations complete!")
	return nil
}
