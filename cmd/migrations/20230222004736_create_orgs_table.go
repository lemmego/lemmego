package main

import (
	"database/sql"
	// "fmt"

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
	q1 := migration.NewSchema().Create("orgs", func(t *migration.Table) error {
		t.AddColumn("id").Type("bigserial").Primary()
		t.AddColumn("name").Type("varchar(255)")
		t.AddColumn("subdomain").Type("varchar(255)")
		t.AddColumn("email").Type("varchar(255)")
		t.AddColumn("created_at").Type("timestamptz(0)").DefaultValue("now()")
		t.AddColumn("updated_at").Type("timestamptz(0)").DefaultValue("now()")
		t.Unique("subdomain", "email")
		return nil
	}).Build()

	// q := schema.Build()
	// fmt.Println("==========start===========")
	// fmt.Println(q)
	// fmt.Println("==========end===========")
	_, err := tx.Exec(q1)
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
