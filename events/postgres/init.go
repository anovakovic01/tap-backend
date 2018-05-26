package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // required for SQL access
	migrate "github.com/rubenv/sql-migrate"
)

// Connect creates a connection to the PostgreSQL instance and applies any
// unapplied database migrations.
func Connect(host, port, name, user, pass string) (*sql.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, name, pass)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err := migrateDB(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrateDB(db *sql.DB) error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "events_1",
				Up: []string{
					`CREATE TABLE events (
						id          BIGSERIAL PRIMARY KEY,
						owner       TEXT NOT NULL,
						title       TEXT NOT NULL,
						description TEXT NOT NULL,
						lat         DECIMAL NOT NULL,
						lon         DECIMAL NOT NULL,
                        start       TIMESTAMP NOT NULL,
						ending      TIMESTAMP NOT NULL
					)`,
				},
				Down: []string{
					"DROP TABLE events",
				},
			},
		},
	}

	_, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	return err
}
