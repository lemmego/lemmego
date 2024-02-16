package main

import (
	"database/sql"

	"github.com/lemmego/migration"
)

func init() {
	migration.GetMigrator().AddMigration(&migration.Migration{
		Version: "20240110230332",
		Up:      mig_20240110230332_create_posts_table_up,
		Down:    mig_20240110230332_create_posts_table_down,
	})
}

func mig_20240110230332_create_posts_table_up(tx *sql.Tx) error {
	schema := migration.NewSchema().Create("posts", func(t *migration.Table) error {
		t.BigIncrements("id").Primary()
		t.UnsignedBigInt("org_id")
		t.String("title", 255)
		t.Text("body").Nullable()
		t.DateTime("created_at").Default("current_timestamp")
		t.DateTime("updated_at").Default("current_timestamp")
		t.ForeignKey("org_id").References("id").On("orgs").OnDelete("cascade")
		return nil
	}).Build()

	if _, err := tx.Exec(schema); err != nil {
		return err
	}

	return nil
}

func mig_20240110230332_create_posts_table_down(tx *sql.Tx) error {
	schema := migration.NewSchema().Drop("posts").Build()
	if _, err := tx.Exec(schema); err != nil {
		return err
	}
	return nil
}
