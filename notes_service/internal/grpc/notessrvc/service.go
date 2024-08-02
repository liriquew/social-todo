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
	SaveNote(context.Context, int64, *models.Note) (int64, error)
	GetNote(context.Context, int64, int64) (*models.Note, error)
	UpdateNote(context.Context, int64, int64, *models.Note) error
	DeleteNote(context.Context, int64, int64) error

	ListUserNotesID(context.Context, int64) ([]int64, error)
	ListUserNotes(context.Context, int64, []int64) ([]*models.Note, error)

	ListUsersNotes(context.Context, []int64, int64, int64) ([]*models.Note, error)
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

func (s *ServiceNotes) CreateNote(ctx context.Context, UID int64, note *notes.Note) (int64, error) {
	const op = "notessrvc.CreateNote"

	log := s.log.With(slog.String("op", op), slog.Int64("uid", UID))
	log.Info("attempting to Create note")

	noteID, err := s.Storage.SaveNote(ctx, UID, models.NoteFromProto(note))
	if err != nil {
		log.Warn("ERROR:", sl.Err(err))
		if errors.Is(err, storage.ErrAlreadyExists) {
			return 0, ErrAlreadyExists
		}

		return 0, err
	}
	return noteID, nil
}

func (s *ServiceNotes) GetNote(ctx context.Context, UID, NID int64) (*notes.Note, error) {
	const op = "notessrvc.GetNote"

	log := s.log.With(slog.String("op", op), slog.Int64("uid", UID), slog.Int64("uid", NID))
	log.Info("attempting to Get note")

	note, err := s.Storage.GetNote(ctx, UID, NID)
	if err != nil {
		log.Warn("ERROR:", sl.Err(err))
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return models.ProtoFromNote(note), nil
}

func (s *ServiceNotes) UpdateNote(ctx context.Context, UID, NID int64, note *notes.Note) error {
	const op = "notessrvc.UpdateNote"

	log := s.log.With(slog.String("op", op), slog.Int64("uid", UID), slog.Int64("uid", NID))
	log.Info("attempting to Update note")

	err := s.Storage.UpdateNote(ctx, UID, NID, models.NoteFromProto(note))
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

func (s *ServiceNotes) DeleteNote(ctx context.Context, UID, NID int64) error {
	const op = "notessrvc.DeleteNote"

	log := s.log.With(slog.String("op", op), slog.Int64("uid", UID), slog.Int64("uid", NID))
	log.Info("attempting to Delete note")

	if err := s.Storage.DeleteNote(ctx, UID, NID); err != nil {
		log.Warn("ERROR:", sl.Err(err))
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func (s *ServiceNotes) ListUserNotesID(ctx context.Context, UID int64) ([]int64, error) {
	const op = "notessrvc.ListUserNotesID"

	log := s.log.With(slog.String("op", op), slog.Int64("uid", UID))
	log.Info("attempting to List user notes IDs note")

	noteIDs, err := s.Storage.ListUserNotesID(ctx, UID)
	if err != nil {
		log.Warn("ERROR", sl.Err(err))
		return nil, err
	}
	if len(noteIDs) == 0 {
		return nil, ErrNotFound
	}

	return noteIDs, nil
}

func (s *ServiceNotes) ListUserNotes(ctx context.Context, UID int64, notesIDs []int64) ([]*notes.NoteListItem, error) {
	const op = "notessrvc.ListUserNotesID"

	log := s.log.With(slog.String("op", op), slog.Int64("uid", UID))
	log.Info("attempting to List user notes note")

	notesList, err := s.Storage.ListUserNotes(ctx, UID, notesIDs)
	if err != nil {
		log.Warn("ERROR", sl.Err(err))
		return nil, err
	}
	if len(notesList) == 0 {
		return nil, ErrNotFound
	}

	notesListRes := make([]*notes.NoteListItem, 0, len(notesList))
	for _, note := range notesList {
		notesListRes = append(notesListRes, &notes.NoteListItem{
			NID:  note.NID,
			Note: models.ProtoFromNote(note),
		})
	}

	return notesListRes, nil
}

func (s *ServiceNotes) ListUsersNotes(ctx context.Context, UIDs []int64, offset, limit int64) ([]*notes.UsersNotesListItem, error) {
	const op = "notessrvc.ListUsersNotes"

	log := s.log.With(slog.String("op", op))
	log.Info("attempting to list users notes", slog.Any("UIDs", UIDs), slog.Int64("OFFSET", offset), slog.Int64("LIMIT", limit))
	notesList, err := s.Storage.ListUsersNotes(ctx, UIDs, offset, limit)
	if err != nil {
		log.Warn("ERROR", sl.Err(err))
		return nil, err
	}

	if len(notesList) == 0 {
		log.Warn("NOT FOUND")
		return nil, ErrNotFound
	}

	notesListRes := make([]*notes.UsersNotesListItem, 0, len(notesList))
	for _, note := range notesList {
		notesListRes = append(notesListRes, &notes.UsersNotesListItem{
			NID:  note.NID,
			UID:  note.UID,
			Note: models.ProtoFromNote(note),
		})
	}

	return notesListRes, nil
}
