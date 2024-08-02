package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/liriquew/social-todo/notes_service/internal/grpc/notessrvc"
	"github.com/liriquew/todoprotos/gen/go/notes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotesService interface {
	CreateNote(context.Context, int64, *notes.Note) (int64, error)
	GetNote(context.Context, int64, int64) (*notes.Note, error)
	UpdateNote(context.Context, int64, int64, *notes.Note) error
	DeleteNote(context.Context, int64, int64) error

	ListUserNotesID(context.Context, int64) ([]int64, error)
	ListUserNotes(context.Context, int64, []int64) ([]*notes.NoteListItem, error)

	ListUsersNotes(context.Context, []int64, int64, int64) ([]*notes.UsersNotesListItem, error)
}

type serverAPI struct {
	notes.UnimplementedNotesServer
	api NotesService
}

func Register(gRPC *grpc.Server, notesService NotesService) {
	notes.RegisterNotesServer(gRPC, &serverAPI{api: notesService})
}

var (
	ErrUID           = fmt.Errorf("empty UID")
	ErrNID           = fmt.Errorf("empty Note ID")
	ErrEmptyContent  = fmt.Errorf("empty content field")
	ErrEmptyTitle    = fmt.Errorf("empty title field")
	ErrInvalidTime   = fmt.Errorf("invalid time duration")
	ErrInvalidUpdate = fmt.Errorf("invalid update request")
	ErrCreatedAtTS   = fmt.Errorf("createdAt timestamp mush be nil")
	ErrEmptyUIDList  = fmt.Errorf("empty UID list")
	ErrBadLimit      = fmt.Errorf("bad limit value")
	ErrBadOffset     = fmt.Errorf("bad offset value")

	MinTime = time.Minute * 10
)

func (s *serverAPI) CreateNote(ctx context.Context, req *notes.CreateNoteRequest) (*notes.NoteResponse, error) {
	if err := validateRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	noteID, err := s.api.CreateNote(ctx, req.UID, req.Note)

	if err != nil {
		if errors.Is(err, notessrvc.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal error idk")
	}

	return &notes.NoteResponse{NID: noteID}, err
}

func (s *serverAPI) GetNote(ctx context.Context, req *notes.NoteIDRequest) (*notes.Note, error) {
	if err := validateRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	note, err := s.api.GetNote(ctx, req.UID, req.NID)

	if err != nil {
		if errors.Is(err, notessrvc.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal error idk")
	}

	return note, err
}

func (s *serverAPI) UpdateNote(ctx context.Context, req *notes.UpdateNoteRequest) (*notes.NoteResponse, error) {
	if err := validateRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.api.UpdateNote(ctx, req.UID, req.NID, req.Note)

	if err != nil {
		if errors.Is(err, notessrvc.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, notessrvc.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal error idk")
	}
	return &notes.NoteResponse{NID: req.NID}, err
}

func (s *serverAPI) DeleteNote(ctx context.Context, req *notes.NoteIDRequest) (*notes.NoteResponse, error) {
	fmt.Println("DELETE")
	if err := validateRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.api.DeleteNote(ctx, req.UID, req.NID)

	if err != nil {
		// idk is this nessesary
		if errors.Is(err, notessrvc.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal error idk")
	}

	return &notes.NoteResponse{NID: req.NID}, err
}

func (s *serverAPI) ListUserNotesID(ctx context.Context, req *notes.UserIDRequest) (*notes.NoteIDList, error) {
	if req.UID <= 0 {
		return nil, status.Error(codes.InvalidArgument, ErrUID.Error())
	}

	noteIDs, err := s.api.ListUserNotesID(ctx, req.UID)

	if err != nil {
		if errors.Is(err, notessrvc.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal error idk")
	}

	return &notes.NoteIDList{UID: req.UID, NoteIDs: noteIDs}, nil
}

func (s *serverAPI) ListUserNotes(ctx context.Context, req *notes.NoteIDList) (*notes.NoteList, error) {
	if req.UID <= 0 {
		return nil, status.Error(codes.InvalidArgument, ErrUID.Error())
	}
	if len(req.NoteIDs) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty ids list")
	}

	notesList, err := s.api.ListUserNotes(ctx, req.UID, req.NoteIDs)
	fmt.Println(len(notesList))

	for _, n := range notesList {
		fmt.Println(n)
	}

	if err != nil {
		if errors.Is(err, notessrvc.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal error idk")
	}
	return &notes.NoteList{Notes: notesList}, nil
}

func (s *serverAPI) ListUsersNotes(ctx context.Context, req *notes.UsersNotesRequest) (*notes.UsersNotesList, error) {

	notesList, err := s.api.ListUsersNotes(ctx, req.UIDs, req.Offset, req.Limit)
	if err != nil {
		if errors.Is(err, notessrvc.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal err idk")
	}

	return &notes.UsersNotesList{Notes: notesList}, nil
}

func validateRequest(req interface{}) error {
	switch v := req.(type) {
	case *notes.CreateNoteRequest:
		if v.Note.CreatedAt != nil {
			return ErrCreatedAtTS
		}
		if v.UID <= 0 {
			return ErrUID
		}
		if v.Note.Content == "" {
			return ErrEmptyContent
		}
		if v.Note.Title == "" {
			return ErrEmptyTitle
		}
		if v.Note.Duration.AsDuration() < MinTime {
			return ErrInvalidTime
		}
	case *notes.NoteIDRequest:
		if v.NID <= 0 {
			return ErrNID
		}
		if v.UID <= 0 {
			return ErrUID
		}
	case *notes.UpdateNoteRequest:
		if v.Note.CreatedAt != nil {
			return ErrCreatedAtTS
		}
		if v.NID <= 0 {
			return ErrNID
		}
		if v.UID <= 0 {
			return ErrUID
		}
		if v.Note.Content == "" && v.Note.Title == "" && v.Note.Duration.AsDuration() < MinTime {
			return ErrInvalidUpdate
		}
		if v.Note.Duration.AsDuration() != 0 && v.Note.Duration.AsDuration() < MinTime {
			return ErrInvalidTime
		}
	case *notes.UsersNotesRequest:
		if v.Limit <= 0 {
			return ErrBadLimit
		}
		if v.Offset <= 0 {
			return ErrBadOffset
		}
		if len(v.UIDs) == 0 {
			return ErrEmptyUIDList
		}
		for _, UID := range v.UIDs {
			if UID <= 0 {
				return ErrUID
			}
		}
	}
	return nil
}
