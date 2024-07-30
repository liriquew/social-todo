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

func (s *Storage) SaveNote(ctx context.Context, UID int64, note *models.Note) (int64, error) {
	const op = "postgres.SaveNote"
	stmt, err := s.db.Prepare("INSERT INTO notes (owner_id, title, note, duration, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id")
	if err != nil {
		return 0, err
	}

	var NID int64
	err = stmt.QueryRow(UID, note.Title, note.Content, note.Duration, time.Now()).Scan(&NID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrAlreadyExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return NID, nil
}

func (s *Storage) GetNote(ctx context.Context, UID, noteID int64) (*models.Note, error) {
	const op = "postgres.GetNote"

	stmt, err := s.db.Preparex("SELECT title, note, duration, created_at FROM notes WHERE owner_id=$1 and id=$2")
	if err != nil {
		return nil, err
	}

	var note models.Note
	if err := stmt.Get(&note, UID, noteID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &note, nil
}

func (s *Storage) UpdateNote(ctx context.Context, UID, NID int64, note *models.Note) error {
	const op = "postgres.UpdateNote"
	stmt, err := s.db.Prepare(`
		UPDATE notes 
		SET 
			title=CASE WHEN $1::text<>'' THEN $1 ELSE title END, 
			note=CASE WHEN $2::text<>'' THEN $2 ELSE note END, 
			duration=CASE WHEN $3::bigint<>0 THEN $3 ELSE duration END
		WHERE 
			owner_id=$4 and id=$5 
		RETURNING id;`)
	if err != nil {
		return err
	}

	var id int64
	err = stmt.QueryRow(note.Title, note.Content, note.Duration, UID, NID).Scan(&id)
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

func (s *Storage) DeleteNote(ctx context.Context, UID, NID int64) error {
	const op = "postgres.DeleteNote"
	stmt, err := s.db.Preparex("DELETE FROM notes WHERE owner_id=$1 and id=$2")
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(UID, NID); err != nil {
		// TODO: idk mb check not exists
		// но это идемподентная операция хз
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) ListUserNotesID(ctx context.Context, UID int64) ([]int64, error) {
	const op = "postgres.ListUserNotedID"
	stmt, err := s.db.Preparex("SELECT id FROM notes WHERE owner_id=$1")
	if err != nil {
		return nil, err
	}

	var NIDs []int64
	err = stmt.Select(&NIDs, UID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return NIDs, err
}

func (s *Storage) ListUserNotes(ctx context.Context, UID int64, NIDs []int64) ([]*models.Note, error) {
	const op = "postgres.ListUserNoted"

	query, args, err := sqlx.In("SELECT id, title, note, duration FROM notes WHERE owner_id=? and id IN(?)", UID, NIDs)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	query = s.db.Rebind(query)

	var notes []*models.Note
	err = s.db.SelectContext(ctx, &notes, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return notes, err
}
