// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // required for SQL access
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
)

// Config defines the options that are used when connecting to a PostgreSQL instance
type Config struct {
	Host        string
	Port        string
	User        string
	Pass        string
	Name        string
	SSLMode     string
	SSLCert     string
	SSLKey      string
	SSLRootCert string
}

// Connect creates a connection to the PostgreSQL instance and applies any
// unapplied database migrations. A non-nil error is returned to indicate
// failure.
func Connect(cfg Config) (*sqlx.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s", cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Pass, cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootCert)

	db, err := sqlx.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	if err := migrateDB(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrateDB(db *sqlx.DB) error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "messages_1",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS messages (
						subtopic      VARCHAR(254),
						publisher     UUID,
						protocol      TEXT,
						name          TEXT,
						unit          TEXT,
						value         FLOAT,
						string_value  TEXT,
						bool_value    BOOL,
						data_value    BYTEA,
						sum           FLOAT,
						time          FLOAT,
						update_time   FLOAT,
						PRIMARY KEY   (time, publisher, subtopic, name)
					)`,
				},
				Down: []string{
					"DROP TABLE messages",
					"DROP TABLE json",
				},
			},
			{
				Id: "messages_2",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS json (
						created       BIGINT,
						subtopic      VARCHAR(254),
						publisher     VARCHAR(254),
						protocol      TEXT,
						payload       JSONB
					)`,
				},
				Down: []string{
					"DROP TABLE json",
				},
			},
			{
				Id: "messages_3",
				Up: []string{
					`ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_pkey`,
				},
			},
			{
				Id: "messages_4",
				Up: []string{
					`ALTER TABLE messages ALTER COLUMN time TYPE BIGINT USING CAST(time AS BIGINT);`,
				},
			},
		},
	}

	_, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	return err
}
