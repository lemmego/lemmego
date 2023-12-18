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
	return nil
}

func mig_20231128193645_create_users_table_down(tx *sql.Tx) error {
	return nil
}
