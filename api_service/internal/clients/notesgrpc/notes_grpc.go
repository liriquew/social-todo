package notes_grpc

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/liriquew/social-todo/api_service/internal/lib/config"
	"github.com/liriquew/social-todo/api_service/internal/models"
	"github.com/liriquew/todoprotos/gen/go/notes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
)

type Client struct {
	api notes.NotesClient
	log *slog.Logger
}

func New(log *slog.Logger, cfg config.ServiceConfig) (*Client, error) {
	const op = "notes_grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(cfg.Retries)),
		grpcretry.WithPerRetryTimeout(cfg.Timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.NewClient(cfg.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{
		api: notes.NewNotesClient(cc),
		log: log,
	}, nil
}

func InterceptorLogger(log *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		log.Log(ctx, slog.Level(level), msg, fields...)
	})
}

var (
	ErrNotFound        = fmt.Errorf("note not found")
	ErrAlreadyExists   = fmt.Errorf("note already exists")
	ErrInvalidArgument = fmt.Errorf("invalid argument")
)

func (c *Client) Create(ctx context.Context, uid int64, note *models.Note) (int64, error) {
	const op = "notes_grpc.Create"

	resp, err := c.api.CreateNote(ctx, &notes.Note{
		Uid:     uid,
		Title:   note.Title,
		Content: note.Note,
	})
	if err != nil {
		// TODO: check status and return this level err
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return 0, ErrInvalidArgument
			case codes.AlreadyExists:
				return 0, ErrAlreadyExists
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return resp.NoteID, nil
}

func (c *Client) Get(ctx context.Context, UID, noteID int64) (*models.Note, error) {
	const op = "notes_grpc.Get"

	resp, err := c.api.GetNoteByID(ctx, &notes.NoteMeta{
		UID:    UID,
		NoteID: noteID,
	})
	if err != nil {
		// TODO: check status and return this level err
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &models.Note{
		Title: resp.Title,
		Note:  resp.Content,
	}, nil
}

func (c *Client) Update(ctx context.Context, UID, noteID int64, note *models.Note) error {
	const op = "notes_grpc.Update"

	_, err := c.api.UpdateNoteByID(ctx, &notes.NoteWithID{
		Meta: &notes.NoteMeta{
			NoteID: noteID,
			UID:    UID,
		},
		Title:   note.Title,
		Content: note.Note,
	})
	if err != nil {
		// TODO: not found idk
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return ErrInvalidArgument
			case codes.AlreadyExists:
				return ErrAlreadyExists
			case codes.NotFound:
				return ErrNotFound
			}
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (c *Client) Delete(ctx context.Context, UID, noteID int64) error {
	const op = "notes_grpc.Delete"

	_, err := c.api.DeleteNotebyID(ctx, &notes.NoteMeta{
		UID:    UID,
		NoteID: noteID,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
