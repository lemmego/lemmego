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
	if _, err := tx.Exec(`create table "posts" (
		"id" bigserial,
		"org_id" bigint not null,
		"title" varchar(255) not null,
		"body" text,
		"created_at" timestamptz(0) default current_timestamp not null,
		"updated_at" timestamptz(0) default current_timestamp not null
	);
	alter table "posts" add constraint "posts_org_id_foreign" foreign key ("org_id") references "orgs" ("id") on delete cascade;
	`); err != nil {
		return err
	}

	return nil
}

func mig_20240110230332_create_posts_table_down(tx *sql.Tx) error {
	if _, err := tx.Exec(`drop table if exists "posts" cascade;`); err != nil {
		return err
	}

	return nil
}
