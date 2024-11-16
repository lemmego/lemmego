package migrations

import (
	"database/sql"
	"github.com/lemmego/migration"
)

func init() {
	migration.GetMigrator().AddMigration(&migration.Migration{
		Version: "20241117024753",
		Up:      mig_20241117024753_create_sessions_table_up,
		Down:    mig_20241117024753_create_sessions_table_down,
	})
}

func mig_20241117024753_create_sessions_table_up(tx *sql.Tx) error {
	schema := migration.Create("sessions", func(t *migration.Table) {
		t.String("token", 255).Primary()
		t.Text("data")
		t.Timestamp("expiry", 6)
	}).Build()

	if _, err := tx.Exec(schema); err != nil {
		return err
	}

	return nil
}

func mig_20241117024753_create_sessions_table_down(tx *sql.Tx) error {
	schema := migration.Drop("sessions").Build()
	if _, err := tx.Exec(schema); err != nil {
		return err
	}
	return nil
}
