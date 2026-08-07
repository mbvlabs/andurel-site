package config

import (
	"fmt"
	"net"
	"net/url"

	"github.com/caarlos0/env/v11"
)

type Database struct {
	Port         string `env:"DB_PORT"`
	Host         string `env:"DB_HOST"`
	Name         string `env:"DB_NAME"`
	User         string `env:"DB_USER"`
	Password     string `env:"DB_PASSWORD"`
	DatabaseKind string `env:"DB_KIND"`
	SslMode      string `env:"DB_SSL_MODE"`
}

func (d Database) GetDatabaseURL() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		d.DatabaseKind, d.User, d.Password, d.Host, d.Port,
		d.Name, d.SslMode,
	)
}

// PostgresURL builds a postgres connection string from individual parts.
func PostgresURL(host, port, name, user, password, sslMode string) string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(host, port),
		Path:   "/" + name,
	}

	q := u.Query()
	q.Set("sslmode", sslMode)
	u.RawQuery = q.Encode()

	return u.String()
}

func newDatabaseConfig() Database {
	dataCfg := Database{}

	if err := env.Parse(&dataCfg); err != nil {
		panic(err)
	}

	return dataCfg
}
