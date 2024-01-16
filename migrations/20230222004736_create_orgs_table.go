package main

import (
	"database/sql"

	"github.com/lemmego/migration"
)

func init() {
	migration.GetMigrator().AddMigration(&migration.Migration{
		Version: "20230222004736",
		Up:      mig_20230222004736_create_orgs_table_up,
		Down:    mig_20230222004736_create_orgs_table_down,
	})
}

func mig_20230222004736_create_orgs_table_up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create table "orgs" (
			"id" bigserial primary key,
			"name" varchar(255) not null,
			"subdomain" varchar(255) not null,
			"email" varchar(255) not null,
			"created_at" timestamptz(0) not null,
			"updated_at" timestamptz(0) not null
		);
		alter table "orgs" add constraint "orgs_email_unique" unique ("email");
	`)
	if err != nil {
		return err
	}
	return nil
}

func mig_20230222004736_create_orgs_table_down(tx *sql.Tx) error {
	_, err := tx.Exec(`drop table if exists "orgs" cascade;`)
	if err != nil {
		return err
	}
	return nil
}
