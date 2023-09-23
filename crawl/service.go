package crawl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

type Service struct {
	db *sql.DB
}

func NewService(
	ctx context.Context,
	filename string,
) (*Service, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, fmt.Errorf("could not open the database: %w", err)
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS usernames (
			id           INTEGER PRIMARY KEY AUTOINCREMENT,
			username     TEXT NOT NULL,
			rank         REAL,
			processed_at TIMESTAMP,
			typename     TEXT CHECK ( typename IN ('User', 'Organization')),
			created_at   TIMESTAMP DEFAULT (datetime('now','localtime'))
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_usernames_username ON usernames(username);
		CREATE TABLE IF NOT EXISTS followers (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			from_id    INT,
			to_id      INT,
			created_at TIMESTAMP DEFAULT (datetime('now','localtime'))
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_followers ON followers(from_id, to_id);

		pragma journal_mode = WAL;
		pragma synchronous = normal;
		pragma temp_store = memory;
		pragma mmap_size = 30000000000;
		pragma auto_vacuum = incremental; -- once on first DB create
		pragma incremental_vacuum; -- regularily
	`)
	if err != nil {
		return nil, fmt.Errorf("could not create the schema: %w", err)
	}

	return &Service{
		db: db,
	}, nil
}

func (s *Service) SetUsername(
	ctx context.Context,
	username string,
) error {
	_, err := s.db.ExecContext(ctx, `INSERT OR IGNORE INTO usernames(username, typename) VALUES (?, "User")`, username)
	if err != nil {
		slog.Error("adding username failed",
			slog.String("username", username),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("could not insert username: %w", err)
	}

	slog.Info("added username", slog.String("username", username))

	return nil
}

func (s *Service) NextUsername(
	ctx context.Context,
) (string, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT
			username
		FROM usernames
		WHERE
			processed_at IS NULL AND
			typename = 'User'
		LIMIT 1;
	`)
	if row.Err() != nil {
		return "", fmt.Errorf("could not query find next user: %w", row.Err())
	}

	var username string
	err := row.Scan(&username)

	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("there are no more usernames to process")
	}

	if err != nil {
		return "", fmt.Errorf("could not read next user: %w", err)
	}

	return username, nil
}

func (s *Service) SetFollower(
	ctx context.Context,
	from string,
	to string,
) error {
	_ = s.SetUsername(ctx, from)
	_ = s.SetUsername(ctx, to)

	_, err := s.db.ExecContext(ctx, `
	INSERT OR IGNORE INTO followers (from_id, to_id)
		SELECT
  		(SELECT id FROM usernames WHERE username = @from) as from_id,
  		(SELECT id FROM usernames WHERE username = @to) as to_id;

	`,
		sql.Named("from", from),
		sql.Named("to", to),
	)

	if err != nil {
		slog.Error(
			"relationship already exists",
			slog.String("from", from),
			slog.String("to", to),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("could not insert follower: %w", err)
	}

	slog.Info(
		"relationship created",
		slog.String("from", from),
		slog.String("to", to),
	)

	return nil
}

func (s *Service) SetProcessed(
	ctx context.Context,
	username string,
) error {
	_, err := s.db.ExecContext(ctx, `UPDATE usernames SET processed_at = datetime('now','localtime') WHERE username = ?`, username)
	if err != nil {
		return fmt.Errorf("could not process %q: %w", username, err)
	}

	return nil
}

func (s *Service) Close() error {
	_, err := s.db.Exec(`
		PRAGMA vacuum;
		PRAGMA optimize;
	`)
	if err != nil {
		return fmt.Errorf("could not optimize database: %w", err)
	}

	err = s.db.Close()
	if err != nil {
		return fmt.Errorf("could not close DB: %w", err)
	}

	return nil
}