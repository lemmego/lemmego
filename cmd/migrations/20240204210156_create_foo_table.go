package main

import (
	"database/sql"

	"github.com/lemmego/migration"
)

func init() {
	migration.GetMigrator().AddMigration(&migration.Migration{
		Version: "20240204210156",
		Up:      mig_20240204210156_create_foo_table_up,
		Down:    mig_20240204210156_create_foo_table_down,
	})
}

func mig_20240204210156_create_foo_table_up(tx *sql.Tx) error {
	q := migration.NewSchema().
		Create("foo", func(t *migration.Table) error {
			t.AddColumn("id").Type("bigserial").Primary()
			return nil
		}).Build()

	_, err := tx.Exec(q)
	return err
}

func mig_20240204210156_create_foo_table_down(tx *sql.Tx) error {
	q := migration.NewSchema().
		Drop("foo").Build()

	_, err := tx.Exec(q)
	return err
}
