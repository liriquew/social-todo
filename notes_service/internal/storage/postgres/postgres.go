package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lib/pq"
	"github.com/liriquew/social-todo/notes_service/internal/lib/config"
	"github.com/liriquew/social-todo/notes_service/internal/models"
	"github.com/liriquew/social-todo/notes_service/internal/storage"
)

type Storage struct {
	db *sqlx.DB
}

func New(cfg config.Config) (*Storage, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", // дада ssl всегда выкл
		cfg.PostgresCfg.Username,
		cfg.PostgresCfg.Password,
		cfg.PostgresCfg.Port,
		cfg.PostgresCfg.DBName,
	)

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) SaveNote(ctx context.Context, note *models.Note) (int64, error) {
	const op = "postgres.SaveNote"
	stmt, err := s.db.Prepare("INSERT INTO notes (owner_id, title, note, created_at) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return 0, err
	}

	var id int64
	note.CreatedTime = time.Now()
	err = stmt.QueryRow(note.UID, note.Title, note.Content, note.CreatedTime).Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrAlreadyExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetNote(ctx context.Context, UID, noteID int64) (*models.Note, error) {
	const op = "postgres.GetNote"
	// возможно методом денотационной семантики можно доказать,
	// что owner_id брать необязательно (вот это слово рил хз как пишется)
	stmt, err := s.db.Preparex("SELECT owner_id, title, note FROM notes WHERE owner_id = $1 and id = $2")
	if err != nil {
		return nil, err
	}

	var note models.Note
	if err := stmt.Get(&note, UID, noteID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}

		fmt.Println(err)

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &note, nil
}

func (s *Storage) UpdateNote(ctx context.Context, note *models.Note) error {
	const op = "postgres.UpdateNote"
	stmt, err := s.db.Preparex("UPDATE notes SET title=$1, note=$2 WHERE id=$3 and owner_id=$4 RETURNING id")
	if err != nil {
		return err
	}

	var id int64
	err = stmt.QueryRow(note.Title, note.Content, note.NoteID, note.UID).Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrAlreadyExists)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteNote(ctx context.Context, UID, noteID int64) error {
	const op = "postgres.DeleteNote"
	stmt, err := s.db.Preparex("DELETE FROM notes WHERE owner_id=$1 and id=$2")
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(UID, noteID); err != nil {
		// TODO: idk mb check not exists
		// но это идемподентная операция хз
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
