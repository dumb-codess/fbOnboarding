package db

import (
	"context"
	"fbOnboarding/ent"

	_ "github.com/mattn/go-sqlite3"
)

var client *ent.Client

type DatabaseConfig struct {
	Dialect    string `envconfig:"DIALECT"`
	ConnString string `envconfig:"CONN_STRING"`
}

func Init(ctx context.Context, cfg DatabaseConfig) error {
	dbclient, err := ent.Open(cfg.Dialect, cfg.ConnString)
	if err != nil {
		return err
	}
	// Run the automatic migration tool to create all schema resources.
	if err := dbclient.Schema.Create(ctx); err != nil {
		return err
	}

	client = dbclient

	return nil
}

func GetDBlient() *ent.Client {
	return client
}
