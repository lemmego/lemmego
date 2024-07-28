package main

import (
	"database/sql"
	"github.com/lemmego/migration"
)

func init() {
	migration.GetMigrator().AddMigration(&migration.Migration{
		Version: "20240726225549",
		Up:      mig_20240726225549_create_oauth_clients_table_up,
		Down:    mig_20240726225549_create_oauth_clients_table_down,
	})
}

func mig_20240726225549_create_oauth_clients_table_up(tx *sql.Tx) error {
	schema := migration.Create("oauth_clients", func(t *migration.Table) {
		t.String("id", 255).Primary()
		t.String("secret", 255)
		t.String("name", 255)
		t.Text("redirect_uri")
		t.DateTime("created_at", 0).Nullable()
		t.DateTime("updated_at", 0).Nullable()
		t.DateTime("deleted_at", 0).Nullable()
	}).Build()

	if _, err := tx.Exec(schema); err != nil {
		return err
	}

	return nil
}

func mig_20240726225549_create_oauth_clients_table_down(tx *sql.Tx) error {
	schema := migration.Drop("oauth_clients").Build()

	if _, err := tx.Exec(schema); err != nil {
		return err
	}

	return nil
}
