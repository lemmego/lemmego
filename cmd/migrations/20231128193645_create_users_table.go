package main

import (
	"database/sql"

	"github.com/lemmego/migration"
)

func init() {
	migration.GetMigrator().AddMigration(&migration.Migration{
		Version: "20231128193645",
		Up:      mig_20231128193645_create_users_table_up,
		Down:    mig_20231128193645_create_users_table_down,
	})
}

func mig_20231128193645_create_users_table_up(tx *sql.Tx) error {
	schema := migration.NewSchema().Create("users", func(t *migration.Table) error {
		t.BigIncrements("id").Primary()
		t.Int("github_id").Unique()
		t.Int("org_id")
		t.String("first_name", 255)
		t.String("last_name", 255)
		t.String("email", 255).Unique()
		t.String("password", 255)
		t.DateTime("created_at").Default("now()")
		t.DateTime("updated_at").Default("now()")
		return nil
	}).Build()
	if _, err := tx.Exec(schema); err != nil {
		return err
	}

	return nil
}

func mig_20231128193645_create_users_table_down(tx *sql.Tx) error {
	schema := migration.NewSchema().Drop("users").Build()
	if _, err := tx.Exec(schema); err != nil {
		return err
	}
	return nil
}
