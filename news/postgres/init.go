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
				Id: "news_1",
				Up: []string{
					`CREATE TABLE news (
						id          BIGSERIAL PRIMARY KEY,
						title       TEXT NOT NULL,
						link        TEXT NOT NULL,
						description TEXT NOT NULL,
                        image_title TEXT NOT NULL,
                        image       TEXT NOT NULL,
                        pub_date    TIMESTAMP UNIQUE NOT NULL
					)`,
				},
				Down: []string{
					"DROP TABLE news",
				},
			},
		},
	}

	_, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	return err
}
