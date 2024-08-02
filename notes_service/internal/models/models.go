package models

import (
	"time"

	"github.com/liriquew/todoprotos/gen/go/notes"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Note struct {
	UID         int64     `db:"owner_id"`
	NID         int64     `db:"id"`
	Title       string    `db:"title"`
	Content     string    `db:"note"`
	CreatedTime time.Time `db:"created_at"`
	Duration    int64     `db:"duration"` // время в минутах
}

func NoteFromProto(n *notes.Note) *Note {
	return &Note{
		Title:       n.Title,
		Content:     n.Content,
		Duration:    n.Duration.AsDuration().Milliseconds(),
		CreatedTime: n.CreatedAt.AsTime(),
	}
}

func ProtoFromNote(n *Note) *notes.Note {
	return &notes.Note{
		Title:     n.Title,
		Content:   n.Content,
		Duration:  durationpb.New(time.Duration(time.Millisecond * time.Duration(n.Duration))),
		CreatedAt: timestamppb.New(n.CreatedTime),
	}
}
