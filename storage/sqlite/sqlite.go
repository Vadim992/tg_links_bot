package sqlite

import (
	"context"
	"database/sql"
	"links_tg-bot/lib/e"
	"links_tg-bot/storage"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewStorageDB(path string) (*Storage, error) {

	db, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, e.Wrap("can't open DB", err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap("can't connect to DB", err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := "INSERT INTO `pages` (`url`, `user_name`) VALUES (?,?)"

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName); err != nil {
		return e.Wrap("can't save data in DB", err)
	}

	return nil
}

func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	q := "SELECT `url` FROM `pages` WHERE `user_name` = ? ORDER BY RANDOM() LIMIT 1"

	var url string
	if err := s.db.QueryRowContext(ctx, q, userName).Scan(&url); err != nil {

		switch err {
		case sql.ErrNoRows:
			return nil, storage.ErrorNoSavedPages
		default:
			return nil, e.Wrap("can't pick random page", err)
		}

	}

	return &storage.Page{
		URL:      url,
		UserName: userName,
	}, nil
}

func (s *Storage) Remove(ctx context.Context, page *storage.Page) error {
	q := "DELETE FROM `pages` WHERE `url` =  ? AND `user_name` = ?"

	if _, err := s.db.ExecContext(ctx, q, page.URL, page.UserName); err != nil {
		return e.Wrap("can't remove data from db", err)
	}
	return nil
}

func (s *Storage) IsExist(ctx context.Context, page *storage.Page) (bool, error) {
	q := "SELECT COUNT(*) FROM `pages` WHERE `url` =  ? AND `user_name` = ?"

	var count int
	if err := s.db.QueryRowContext(ctx, q, page.URL, page.UserName).Scan(&count); err != nil {
		return false, e.Wrap("can't check if page exists", err)
	}
	return count > 0, nil
}

func (s *Storage) InitDB(ctx context.Context) error {

	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`

	if _, err := s.db.ExecContext(ctx, q); err != nil {
		return e.Wrap("can't create table", err)
	}

	return nil
}
