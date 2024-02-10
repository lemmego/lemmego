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
	if _, err := tx.Exec(`create table "users" (
		"id" bigserial primary key,
		"github_id" int8 unique,
		"org_id" int8 not null,
		"first_name" varchar(255) not null,
		"last_name" varchar(255),
		"email" varchar(255) unique not null,
		"password" varchar(255) not null,
		"created_at" timestamptz(0) default current_timestamp not null,
		"updated_at" timestamptz(0) default current_timestamp not null
	);`); err != nil {
		return err
	}

	return nil
}

func mig_20231128193645_create_users_table_down(tx *sql.Tx) error {
	if _, err := tx.Exec(`drop table "users";`); err != nil {
		return err
	}
	return nil
}
