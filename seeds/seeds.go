package seeds

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"andurel-site/config"
	"andurel-site/models"
	"andurel-site/models/factories"

	"github.com/mbvlabs/andurel/pkg/storage"
)

const Default = "development"

const (
	productionAdminEmail = "morten@mbvlabs.com"
)

var productionAdminPassword = os.Getenv("ADMIN_PASSWORD")

type Runner func(context.Context, storage.Connection) error

var Registry = map[string]Runner{
	"default":     Development,
	"development": Development,
	"test":        Test,
	"production":  Production,
}

func Names() []string {
	names := make([]string, 0, len(Registry))
	for name := range Registry {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

func Run(ctx context.Context, db storage.Connection, name string) error {
	if name == "" {
		name = Default
	}

	runner, ok := Registry[name]
	if !ok {
		return fmt.Errorf("unknown seed %q (available: %s)", name, strings.Join(Names(), ", "))
	}

	return runner(ctx, db)
}

func Development(ctx context.Context, db storage.Connection) error {
	admin, err := factories.CreateUser(ctx, db,
		factories.WithEmail("admin@example.com"),
		factories.WithIsAdmin(true),
		factories.WithValidatedEmail(),
	)
	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}
	fmt.Printf("Created admin user: %s\n", admin.Email)

	user, err := factories.CreateUser(ctx, db,
		factories.WithEmail("user@example.com"),
		factories.WithValidatedEmail(),
	)
	if err != nil {
		return fmt.Errorf("failed to create regular user: %w", err)
	}
	fmt.Printf("Created regular user: %s\n", user.Email)

	// Add more seeds here using factories:
	//
	// // Create 10 additional users with random emails
	// users, err := factories.CreateUsers(ctx, db, 10)
	// if err != nil {
	// 	return fmt.Errorf("failed to create users: %w", err)
	// }
	// fmt.Printf("Created %d additional users\n", len(users))

	return nil
}

func Test(ctx context.Context, db storage.Connection) error {
	_, err := factories.CreateUser(ctx, db,
		factories.WithEmail("test@example.com"),
		factories.WithValidatedEmail(),
	)
	if err != nil {
		return fmt.Errorf("failed to create test user: %w", err)
	}

	return nil
}

// Production inserts the console admin once. Re-running is a no-op if that
// email already exists; it does not rotate the password.
func Production(ctx context.Context, db storage.Connection) error {
	if productionAdminPassword == "" || productionAdminPassword == "CHANGE_ME" {
		return fmt.Errorf(
			"set productionAdminPassword in seeds/seeds.go before running the production seed",
		)
	}

	auth, err := config.NewAuth()
	if err != nil {
		return err
	}

	users := models.NewUsers(db)
	existing, err := users.FindByEmail(ctx, productionAdminEmail)
	if err == nil {
		fmt.Printf("Production admin already present: %s\n", existing.Email)
		return nil
	}
	if !errors.Is(err, models.ErrNotFound) {
		return fmt.Errorf("lookup production admin: %w", err)
	}

	hashed, err := models.HashPassword(productionAdminPassword, auth.Pepper)
	if err != nil {
		return fmt.Errorf("hash production admin password: %w", err)
	}

	admin, err := factories.CreateUser(ctx, db,
		factories.WithEmail(productionAdminEmail),
		factories.WithIsAdmin(true),
		factories.WithValidatedEmail(),
		factories.WithPassword([]byte(hashed)),
	)
	if err != nil {
		return fmt.Errorf("create production admin: %w", err)
	}

	fmt.Printf("Created production admin user: %s\n", admin.Email)
	return nil
}
