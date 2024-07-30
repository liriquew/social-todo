package models

import (
	"time"

	"github.com/liriquew/todoprotos/gen/go/notes"
	"google.golang.org/protobuf/types/known/durationpb"
)

type Note struct {
	NID       int64     `json:"id,omitempty"`
	Title     string    `json:"title,omitempty"`
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Duration  int64     `json:"duration,omitempty"`
}

func NoteFromProto(n *notes.Note, NID ...int64) *Note {
	var noteID int64
	if len(NID) != 0 {
		noteID = NID[0]
	}
	return &Note{
		NID:       noteID,
		Title:     n.Title,
		Content:   n.Content,
		Duration:  int64(n.Duration.AsDuration().Minutes()),
		CreatedAt: n.CreatedAt.AsTime(),
	}
}

func ProtoFromNote(n *Note) *notes.Note {
	return &notes.Note{
		Title:   n.Title,
		Content: n.Content,
		Duration: &durationpb.Duration{
			Seconds: n.Duration * 60,
		},
	}
}
