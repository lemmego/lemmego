package migrations

import (
	"database/sql"
	"github.com/lemmego/migration"
)

func init() {
	migration.GetMigrator().AddMigration(&migration.Migration{
		Version: "20240814152548",
		Up:      mig_20240814152548_create_users_table_up,
		Down:    mig_20240814152548_create_users_table_down,
	})
}

func mig_20240814152548_create_users_table_up(tx *sql.Tx) error {
	schema := migration.Create("users", func(t *migration.Table) {
		t.BigIncrements("id")
		t.String("email", 255).Unique()
		t.Text("password")
		t.String("first_name", 255)
		t.String("last_name", 255)
		t.String("username", 255).Unique()
		t.Text("bio").Nullable()
		t.String("phone", 255)
		t.String("avatar", 255)
		t.ForeignID("org_id").Constrained()
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

func mig_20240814152548_create_users_table_down(tx *sql.Tx) error {
	schema := migration.Drop("users").Build()

	if _, err := tx.Exec(schema); err != nil {
		return err
	}

	return nil
}
