package models

type Note struct {
	Title string `json:"title,omitempty"`
	Note  string `json:"note,omitempty"`
}

type NoteMeta struct {
	UID    int64 `json:"uid,omitempty"`
	NoteID int64 `json:"note_id,omitempty"`
}
