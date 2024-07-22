package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/liriquew/social-todo/notes_service/internal/grpc/notessrvc"
	"github.com/liriquew/todoprotos/gen/go/notes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotesService interface {
	CreateNote(context.Context, *notes.Note) (*notes.NoteMeta, error)
	GetNote(context.Context, *notes.NoteMeta) (*notes.Note, error)
	UpdateNote(context.Context, *notes.NoteWithID) error
	DeleteNote(context.Context, *notes.NoteMeta) error
}

type serverAPI struct {
	notes.UnimplementedNotesServer
	api NotesService
}

func Register(gRPC *grpc.Server, notesService NotesService) {
	notes.RegisterNotesServer(gRPC, &serverAPI{api: notesService})
}

var (
	ErrUID          = fmt.Errorf("empty UID")
	ErrNoteID       = fmt.Errorf("empty Note ID")
	ErrEmptyContent = fmt.Errorf("empty content field")
	ErrEmptyTitle   = fmt.Errorf("empty title field")
)

func (s *serverAPI) CreateNote(ctx context.Context, req *notes.Note) (*notes.NoteMeta, error) {
	if err := validateRequest(req); err != nil {
		if errors.Is(err, ErrUID) {
			return nil, status.Error(codes.Internal, "miss uid from jwt token")
		}
		if errors.Is(err, ErrEmptyContent) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, ErrEmptyTitle) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	noteMeta, err := s.api.CreateNote(ctx, req)

	if err != nil {
		if errors.Is(err, notessrvc.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal error idk")
	}

	return noteMeta, err
}

func (s *serverAPI) GetNoteByID(ctx context.Context, req *notes.NoteMeta) (*notes.Note, error) {
	if err := validateRequest(req); err != nil {
		if errors.Is(err, ErrUID) {
			return nil, status.Error(codes.Internal, "miss uid from jwt token")
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	note, err := s.api.GetNote(ctx, req)

	// TODO: check err
	if err != nil {
		if errors.Is(err, notessrvc.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error idk")
	}

	return note, err
}

func (s *serverAPI) UpdateNoteByID(ctx context.Context, req *notes.NoteWithID) (*notes.NoteMeta, error) {
	if err := validateRequest(req); err != nil {
		if errors.Is(err, ErrUID) {
			return nil, status.Error(codes.Internal, "miss uid from jwt token")
		}
		if errors.Is(err, ErrEmptyContent) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, ErrEmptyTitle) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.api.UpdateNote(ctx, req)

	// TODO: check err
	if err != nil {
		if errors.Is(err, notessrvc.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, notessrvc.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal error idk")
	}
	return req.Meta, err
}

func (s *serverAPI) DeleteNotebyID(ctx context.Context, req *notes.NoteMeta) (*notes.NoteMeta, error) {
	fmt.Println("DELETE")
	if err := validateRequest(req); err != nil {
		if errors.Is(err, ErrUID) {
			return nil, status.Error(codes.Internal, "miss uid from jwt token")
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.api.DeleteNote(ctx, req)

	// TODO: check err
	if err != nil {
		// idk is this nessesary
		if errors.Is(err, notessrvc.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error idk")
	}

	return req, err
}

func validateRequest(req interface{}) error {
	switch v := req.(type) {
	case *notes.Note:
		if v.Uid == 0 {
			return ErrUID
		}
		if v.Content == "" {
			return ErrEmptyContent
		}
		if v.Title == "" {
			return ErrEmptyTitle
		}
	case *notes.NoteMeta:
		if v.NoteID == 0 {
			return ErrNoteID
		}
		if v.UID == 0 {
			return ErrUID
		}
	case *notes.NoteWithID:
		if v.Meta.NoteID == 0 {
			return ErrNoteID
		}
		if v.Meta.UID == 0 {
			return ErrUID
		}
		if v.Content == "" {
			return ErrEmptyContent
		}
		if v.Title == "" {
			return ErrEmptyTitle
		}
	}
	return nil
}
