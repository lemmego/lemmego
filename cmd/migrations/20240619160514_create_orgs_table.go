package main

import (
  "database/sql"
  "github.com/lemmego/migration"
)

func init() {
  migration.GetMigrator().AddMigration(&migration.Migration{
    Version: "20240619160514",
    Up:      mig_20240619160514_create_orgs_table_up,
    Down:    mig_20240619160514_create_orgs_table_down,
  })
}

func mig_20240619160514_create_orgs_table_up(tx *sql.Tx) error {
  schema := migration.Create("orgs", func(t *migration.Table) {
  	t.BigIncrements("id").Primary()
	t.String("org_username", 255).Unique()
	t.String("org_name", 255)
	t.String("email", 255).Unique()
	t.DateTime("created_at", 0).Nullable()
	t.DateTime("updated_at", 0).Nullable()
	t.DateTime("deleted_at", 0).Nullable()
  }).Build()

  if _, err := tx.Exec(schema); err != nil {
    return err
  }

  return nil
}

func mig_20240619160514_create_orgs_table_down(tx *sql.Tx) error {
  schema := migration.Drop("orgs").Build()

  if _, err := tx.Exec(schema); err != nil {
    return err
  }

  return nil
}
