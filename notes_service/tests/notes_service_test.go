package tests

import (
	"testing"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	"github.com/liriquew/social-todo/notes_service/tests/suite"
	"github.com/liriquew/todoprotos/gen/go/notes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	UID1 int64 = 100
	UID2 int64 = 101
)

func TestCreateRead_Read_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	note := gofakeit.Book()

	respCreate, err := st.NoteClient.CreateNote(ctx, &notes.Note{
		Uid:     UID1,
		Title:   note.Title,
		Content: note.Author,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respCreate.NoteID)
	assert.Equal(t, UID1, respCreate.UID)

	respGet, err := st.NoteClient.GetNoteByID(ctx, &notes.NoteMeta{
		UID:    UID1,
		NoteID: respCreate.NoteID,
	})

	require.NoError(t, err)
	assert.Equal(t, UID1, respGet.Uid)
	assert.Equal(t, note.Title, respGet.Title)
	assert.Equal(t, note.Author, respGet.Content)
}

func TestCreate_OneTitleFewUID(t *testing.T) {
	ctx, st := suite.New(t)

	note := gofakeit.Book()

	respCreate1, err := st.NoteClient.CreateNote(ctx, &notes.Note{
		Uid:     UID1,
		Title:   note.Title,
		Content: note.Author,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respCreate1.NoteID)
	assert.Equal(t, UID1, respCreate1.UID)

	respCreate2, err := st.NoteClient.CreateNote(ctx, &notes.Note{
		Uid:     UID2,
		Title:   note.Title,
		Content: note.Author,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respCreate2.NoteID)
	assert.Equal(t, UID2, respCreate2.UID)

	respGet1, err := st.NoteClient.GetNoteByID(ctx, &notes.NoteMeta{
		UID:    UID1,
		NoteID: respCreate1.NoteID,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respGet1.Content)
	assert.NotEmpty(t, respGet1.Title)
	assert.Equal(t, respGet1.Uid, UID1)

	respGet2, err := st.NoteClient.GetNoteByID(ctx, &notes.NoteMeta{
		UID:    UID2,
		NoteID: respCreate2.NoteID,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respGet2.Content)
	assert.NotEmpty(t, respGet2.Title)
	assert.Equal(t, respGet2.Uid, UID2)

	assert.Equal(t, respGet1.Title, respGet2.Title)
}

func TestCreateCreate_DuplecateCreate(t *testing.T) {
	ctx, st := suite.New(t)

	note := gofakeit.Book()

	respCreate, err := st.NoteClient.CreateNote(ctx, &notes.Note{
		Uid:     UID1,
		Title:   note.Title,
		Content: note.Author,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respCreate.NoteID)
	assert.Equal(t, UID1, respCreate.UID)

	respCreate, err = st.NoteClient.CreateNote(ctx, &notes.Note{
		Uid:     UID1,
		Title:   note.Title,
		Content: note.Author,
	})

	require.Error(t, err)
	assert.Empty(t, respCreate)
	assert.ErrorContains(t, err, "note with that title already exists")
}

func TestCreate_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name          string
		note          notes.Note
		expectedError string
	}{
		{
			name: "Empty content",
			note: notes.Note{
				Uid:     UID1,
				Title:   "some Title",
				Content: "",
			},
			expectedError: "empty content field",
		},
		{
			name: "Empty title",
			note: notes.Note{
				Uid:     UID1,
				Title:   "",
				Content: "some Content",
			},
			expectedError: "empty title field",
		},
		{
			name: "Both empty content",
			note: notes.Note{
				Uid:     UID1,
				Title:   "",
				Content: "",
			},
			expectedError: "empty content field",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.NoteClient.CreateNote(ctx, &tt.note)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

func TestCRUD_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	note := gofakeit.Book()

	respCreate, err := st.NoteClient.CreateNote(ctx, &notes.Note{
		Uid:     UID1,
		Title:   note.Title,
		Content: note.Author,
	})

	assert.NoError(t, err)
	require.NotEmpty(t, respCreate.NoteID)
	require.Equal(t, UID1, respCreate.UID)

	respGet, err := st.NoteClient.GetNoteByID(ctx, &notes.NoteMeta{
		UID:    UID1,
		NoteID: respCreate.NoteID,
	})

	assert.NoError(t, err)
	require.Equal(t, note.Title, respGet.Title)
	require.Equal(t, note.Author, respGet.Content)

	newNote := gofakeit.Book()

	respUpdate, err := st.NoteClient.UpdateNoteByID(ctx, &notes.NoteWithID{
		Meta: &notes.NoteMeta{
			UID:    UID1,
			NoteID: respCreate.NoteID,
		},
		Title:   newNote.Title,
		Content: newNote.Author,
	})

	assert.NoError(t, err)
	require.Equal(t, respCreate.NoteID, respUpdate.NoteID)
	require.Equal(t, respCreate.UID, respUpdate.UID)
	// методом денотационной семантики: respUpdate.UID == UID1

	respNewGet, err := st.NoteClient.GetNoteByID(ctx, &notes.NoteMeta{
		UID:    UID1,
		NoteID: respCreate.NoteID,
	})

	assert.NoError(t, err)
	require.Equal(t, respNewGet.Title, newNote.Title)
	require.Equal(t, respNewGet.Content, newNote.Author)
	require.Equal(t, respNewGet.Uid, UID1)
}

func TestUpdate_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	note := gofakeit.Book()

	respCreate, err := st.NoteClient.CreateNote(ctx, &notes.Note{
		Uid:     UID1,
		Title:   note.Title,
		Content: note.Author,
	})

	assert.NoError(t, err)
	require.NotEmpty(t, respCreate.NoteID)
	require.Equal(t, respCreate.UID, UID1)

	tests := []struct {
		name          string
		note          notes.Note
		expectedError string
	}{
		{
			name: "Empty content",
			note: notes.Note{
				Uid:     UID1,
				Title:   "some Title",
				Content: "",
			},
			expectedError: "empty content field",
		},
		{
			name: "Empty title",
			note: notes.Note{
				Uid:     UID1,
				Title:   "",
				Content: "some Content",
			},
			expectedError: "empty title field",
		},
		{
			name: "Both empty content",
			note: notes.Note{
				Uid:     UID1,
				Title:   "",
				Content: "",
			},
			expectedError: "empty content field",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.NoteClient.UpdateNoteByID(ctx,
				&notes.NoteWithID{
					Meta: &notes.NoteMeta{
						UID:    UID1,
						NoteID: respCreate.NoteID,
					},
					Title:   tt.note.Title,
					Content: tt.note.Content,
				},
			)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

func TestDelete_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	note := gofakeit.Book()

	respCreate, err := st.NoteClient.CreateNote(ctx, &notes.Note{
		Uid:     UID2,
		Title:   note.Title,
		Content: note.Author,
	})

	assert.NoError(t, err)
	require.NotEmpty(t, respCreate.NoteID)
	require.Equal(t, respCreate.UID, UID2)

	respGet, err := st.NoteClient.GetNoteByID(ctx, &notes.NoteMeta{
		UID:    UID2,
		NoteID: respCreate.NoteID,
	})

	assert.NoError(t, err)
	require.Equal(t, note.Title, respGet.Title)
	require.Equal(t, note.Author, respGet.Content)

	respDelete, err := st.NoteClient.DeleteNotebyID(ctx, &notes.NoteMeta{
		UID:    UID2,
		NoteID: respCreate.NoteID,
	})

	assert.NoError(t, err)
	require.Equal(t, respDelete.UID, UID2)
	require.Equal(t, respDelete.NoteID, respCreate.NoteID)

	_, err = st.NoteClient.GetNoteByID(ctx, &notes.NoteMeta{
		UID:    UID2,
		NoteID: respCreate.NoteID,
	})

	assert.Error(t, err)
	require.Contains(t, err.Error(), "note not found")
}
