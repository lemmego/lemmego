package main

import (
  "database/sql"
  "github.com/lemmego/migration"
)

func init() {
  migration.GetMigrator().AddMigration(&migration.Migration{
    Version: "20240410140257",
    Up:      mig_20240410140257_create_posts_table_up,
    Down:    mig_20240410140257_create_posts_table_down,
  })
}

func mig_20240410140257_create_posts_table_up(tx *sql.Tx) error {
  schema := migration.Create("posts", func(t *migration.Table) {
  	t.UnsignedBigInt("id")
	t.String("title", 0)
	t.Text("body").Nullable()
	t.PrimaryKey("id")
  }).Build()

  if _, err := tx.Exec(schema); err != nil {
    return err
  }

  return nil
}

func mig_20240410140257_create_posts_table_down(tx *sql.Tx) error {
  schema := migration.Drop("posts").Build()

  if _, err := tx.Exec(schema); err != nil {
    return err
  }

  return nil
}
