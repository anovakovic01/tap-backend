package postgres

import (
	"database/sql"

	"github.com/anovakovic01/tap-backend/news"
)

var _ news.Repository = (*newsRepository)(nil)

type newsRepository struct {
	db *sql.DB
}

// NewNewsRepositry returns new news repository instance.
func NewNewsRepositry(db *sql.DB) news.Repository {
	return newsRepository{db}
}

func (repo newsRepository) Create(n news.News) (int64, error) {
	q := `INSERT INTO news (title, link, description, image_title, image, pub_date) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var id int64
	err := repo.db.
		QueryRow(q, n.Title, n.Link, n.Description, n.ImageTitle, n.Image, n.PubDate).
		Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo newsRepository) One(id int64) (news.News, error) {
	q := `SELECT title, link, description, image_title, image, pub_date FROM news WHERE id = $1`

	n := news.News{ID: id}
	err := repo.db.
		QueryRow(q, id).
		Scan(&n.Title, &n.Link, &n.Description, &n.ImageTitle, &n.Image, &n.PubDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return news.News{}, news.ErrNotFound
		}
		return news.News{}, err
	}

	return n, nil
}

func (repo newsRepository) All() []news.News {
	q := `SELECT id, title, link, description, image_title, image, pub_date FROM news ORDER BY id`
	items := []news.News{}

	rows, err := repo.db.Query(q)
	if err != nil {
		return []news.News{}
	}
	defer rows.Close()

	for rows.Next() {
		n := news.News{}
		if err := rows.Scan(&n.ID, &n.Title, &n.Link, &n.Description, &n.ImageTitle, &n.Image, &n.PubDate); err != nil {
			return []news.News{}
		}

		items = append(items, n)
	}

	return items
}
