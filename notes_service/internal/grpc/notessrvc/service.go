package notessrvc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/liriquew/social-todo/api_service/pkg/logger/sl"
	"github.com/liriquew/social-todo/notes_service/internal/models"
	"github.com/liriquew/social-todo/notes_service/internal/storage"
	"github.com/liriquew/todoprotos/gen/go/notes"
)

type ServiceNotes struct {
	log     *slog.Logger
	Storage StorageProvider
}

type StorageProvider interface {
	SaveNote(context.Context, *models.Note) (int64, error)
	GetNote(context.Context, int64, int64) (*models.Note, error)
	UpdateNote(context.Context, *models.Note) error
	DeleteNote(context.Context, int64, int64) error
}

var (
	ErrNotFound      = fmt.Errorf("note not found")
	ErrAlreadyExists = fmt.Errorf("note with that title already exists")
)

func New(log *slog.Logger, storage StorageProvider) *ServiceNotes {
	return &ServiceNotes{
		log:     log,
		Storage: storage,
	}
}

func (s *ServiceNotes) CreateNote(ctx context.Context, note *notes.Note) (*notes.NoteMeta, error) {
	const op = "notessrvc.CreateNote"

	log := s.log.With(slog.String("op", op), slog.Int64("uid", note.Uid))
	log.Info("attempting to Create note")

	noteID, err := s.Storage.SaveNote(ctx, models.NoteFromProto(note))
	if err != nil {
		log.Warn("ERROR:", sl.Err(err))
		if errors.Is(err, storage.ErrAlreadyExists) {
			return nil, ErrAlreadyExists
		}

		return nil, err
	}
	return &notes.NoteMeta{UID: note.Uid, NoteID: noteID}, nil
}

func (s *ServiceNotes) GetNote(ctx context.Context, noteMeta *notes.NoteMeta) (*notes.Note, error) {
	const op = "notessrvc.GetNote"

	log := s.log.With(slog.String("op", op), slog.Int64("uid", noteMeta.UID), slog.Int64("uid", noteMeta.NoteID))
	log.Info("attempting to Get note")

	note, err := s.Storage.GetNote(ctx, noteMeta.UID, noteMeta.NoteID)
	if err != nil {
		log.Warn("ERROR:", sl.Err(err))
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &notes.Note{
		Uid:     note.UID,
		Title:   *note.Title,
		Content: *note.Content,
	}, nil
}

func (s *ServiceNotes) UpdateNote(ctx context.Context, noteWithID *notes.NoteWithID) error {
	const op = "notessrvc.UpdateNote"

	log := s.log.With(slog.String("op", op), slog.Int64("uid", noteWithID.Meta.UID), slog.Int64("uid", noteWithID.Meta.NoteID))
	log.Info("attempting to Update note")

	err := s.Storage.UpdateNote(ctx, models.NoteWithIDFromProto(noteWithID))
	if err != nil {
		log.Warn("ERROR:", sl.Err(err))
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}
		if errors.Is(err, storage.ErrAlreadyExists) {
			return ErrAlreadyExists
		}
		return err
	}

	return nil
}

func (s *ServiceNotes) DeleteNote(ctx context.Context, noteMeta *notes.NoteMeta) error {
	const op = "notessrvc.DeleteNote"

	log := s.log.With(slog.String("op", op), slog.Int64("uid", noteMeta.UID), slog.Int64("uid", noteMeta.NoteID))
	log.Info("attempting to Delete note")

	if err := s.Storage.DeleteNote(ctx, noteMeta.UID, noteMeta.NoteID); err != nil {
		log.Warn("ERROR:", sl.Err(err))
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	return nil
}
