package main

import (
  "database/sql"
  "github.com/lemmego/migration"
)

func init() {
  migration.GetMigrator().AddMigration(&migration.Migration{
    Version: "20240720174615",
    Up:      mig_20240720174615_create_users_table_up,
    Down:    mig_20240720174615_create_users_table_down,
  })
}

func mig_20240720174615_create_users_table_up(tx *sql.Tx) error {
  schema := migration.Create("users", func(t *migration.Table) {
    t.BigIncrements("id")
    t.ForeignID("org_id").Constrained()
    t.String("first_name", 255)
    t.String("last_name", 255)
    t.String("logo", 255)
    t.String("email", 255).Unique()
    t.String("password", 255)
    t.DateTime("created_at", 0).Nullable()
    t.DateTime("updated_at", 0).Nullable()
    t.DateTime("deleted_at", 0).Nullable()
    t.PrimaryKey("id, org_id")
  }).Build()

  if _, err := tx.Exec(schema); err != nil {
    return err
  }

  return nil
}

func mig_20240720174615_create_users_table_down(tx *sql.Tx) error {
  schema := migration.Drop("users").Build()

  if _, err := tx.Exec(schema); err != nil {
    return err
  }

  return nil
}
