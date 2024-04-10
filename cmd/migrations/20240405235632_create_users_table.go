package main

import (
	"database/sql"
	"github.com/lemmego/migration"
)

func init() {
	migration.GetMigrator().AddMigration(&migration.Migration{
		Version: "20240405235632",
		Up:      mig_20240405235632_create_users_table_up,
		Down:    mig_20240405235632_create_users_table_down,
	})
}

func mig_20240405235632_create_users_table_up(tx *sql.Tx) error {
	schema := migration.Create("users", func(t *migration.Table) {
		t.UnsignedBigInt("id").Primary()
		t.PrimaryKey("id")
	}).Build()

	if _, err := tx.Exec(schema); err != nil {
		return err
	}

	return nil
}

func mig_20240405235632_create_users_table_down(tx *sql.Tx) error {
	schema := migration.Drop("users").Build()

	if _, err := tx.Exec(schema); err != nil {
		return err
	}

	return nil
}
