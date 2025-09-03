package migrations

import (
	"database/sql"
	"github.com/lemmego/migration"
)

func init() {
	migration.GetMigrator().AddMigration(&migration.Migration{
		Version: "20250903034025",
		Up:      mig_20250903034025_create_users_table_up,
		Down:    mig_20250903034025_create_users_table_down,
	})
}

func mig_20250903034025_create_users_table_up(tx *sql.Tx) error {
	schema := migration.Create("users", func(t *migration.Table) {
		t.UnsignedBigInt("id").Primary()
		t.Text("email").Unique()
		t.Text("name")
		t.Text("password")
		t.DateTime("created_at", 6).Nullable()
		t.DateTime("updated_at", 6).Nullable()
	}).Build()

	if _, err := tx.Exec(schema); err != nil {
		return err
	}

	return nil
}

func mig_20250903034025_create_users_table_down(tx *sql.Tx) error {
	schema := migration.Drop("users").Build()

	if _, err := tx.Exec(schema); err != nil {
		return err
	}

	return nil
}
