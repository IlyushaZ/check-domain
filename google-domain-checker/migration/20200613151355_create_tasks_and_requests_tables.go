package migration

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20200613151355, Down20200613151355)
}

func Up20200613151355(tx *sql.Tx) error {
	_, err := tx.Exec("CREATE TABLE tasks (" +
		"id uuid primary key, " +
		"domain VARCHAR(255) NOT NULL, " +
		"country VARCHAR(255) NOT NULL," +
		"created_at TIMESTAMP NOT NULL, " +
		"process_at TIMESTAMP NOT NULL" +
		");")

	_, err = tx.Exec("CREATE TABLE requests (" +
		"id uuid primary key, " +
		"task_id uuid NOT NULL, " +
		"text varchar(255) NOT NULL, " +
		"FOREIGN KEY (task_id) REFERENCES tasks (id) ON DELETE CASCADE" +
		");")

	return err
}

func Down20200613151355(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE tasks, requests;")
	return err
}
