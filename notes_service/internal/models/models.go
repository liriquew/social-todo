package models

import (
	"time"

	"github.com/liriquew/todoprotos/gen/go/notes"
)

type Note struct {
	UID         int64     `db:"owner_id"`
	NoteID      int64     `db:"id"`
	Title       *string   `db:"title"`
	Content     *string   `db:"note"`
	CreatedTime time.Time `db:"created_at"`
}

func NoteFromProto(n *notes.Note) *Note {
	return &Note{
		UID:     n.Uid,
		Title:   &n.Title,
		Content: &n.Content,
	}
}

func NoteWithIDFromProto(n *notes.NoteWithID) *Note {
	return &Note{
		UID:     n.Meta.UID,
		NoteID:  n.Meta.NoteID,
		Title:   &n.Title,
		Content: &n.Content,
	}
}
