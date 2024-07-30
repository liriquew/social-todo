package tests

import (
	"testing"
	"time"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	"github.com/liriquew/social-todo/notes_service/tests/suite"
	"github.com/liriquew/todoprotos/gen/go/notes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	UID1 int64 = 100
	UID2 int64 = 101
)

func TestCreateRead_Read_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	note := &notes.Note{
		Title:    gofakeit.Name(),
		Content:  gofakeit.HackerPhrase(),
		Duration: durationpb.New(time.Minute * 10),
	}

	respCreate, err := st.NoteClient.CreateNote(ctx, &notes.CreateNoteRequest{
		UID:  UID1,
		Note: note,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respCreate.NID)

	respGet, err := st.NoteClient.GetNote(ctx, &notes.NoteIDRequest{
		UID: UID1,
		NID: respCreate.NID,
	})

	require.NoError(t, err)
	assert.Equal(t, note.Title, respGet.Title)
	assert.Equal(t, note.Content, respGet.Content)
	assert.Equal(t, note.Duration.AsDuration(), respGet.Duration.AsDuration())
}

func TestCreate_OneTitleFewUID(t *testing.T) {
	ctx, st := suite.New(t)

	note := &notes.Note{
		Title:    gofakeit.Name(),
		Content:  gofakeit.HackerPhrase(),
		Duration: durationpb.New(time.Minute * 10),
	}

	respCreate1, err := st.NoteClient.CreateNote(ctx, &notes.CreateNoteRequest{
		UID:  UID1,
		Note: note,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respCreate1.NID)

	respCreate2, err := st.NoteClient.CreateNote(ctx, &notes.CreateNoteRequest{
		UID:  UID2,
		Note: note,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respCreate2.NID)

	respGet1, err := st.NoteClient.GetNote(ctx, &notes.NoteIDRequest{
		UID: UID1,
		NID: respCreate1.NID,
	})

	require.NoError(t, err)
	assert.Equal(t, note.Content, respGet1.Content)
	assert.Equal(t, note.Title, respGet1.Title)
	assert.Equal(t, note.Duration.AsDuration(), respGet1.Duration.AsDuration())

	respGet2, err := st.NoteClient.GetNote(ctx, &notes.NoteIDRequest{
		UID: UID2,
		NID: respCreate2.NID,
	})

	require.NoError(t, err)
	assert.Equal(t, note.Content, respGet1.Content)
	assert.Equal(t, note.Title, respGet1.Title)
	assert.Equal(t, note.Duration.AsDuration(), respGet1.Duration.AsDuration())

	assert.Equal(t, respGet1.Title, respGet2.Title)
}

func TestCreateCreate_DuplecateCreate(t *testing.T) {
	ctx, st := suite.New(t)

	note := &notes.Note{
		Title:    gofakeit.Name(),
		Content:  gofakeit.HackerPhrase(),
		Duration: durationpb.New(time.Minute * 10),
	}

	respCreate, err := st.NoteClient.CreateNote(ctx, &notes.CreateNoteRequest{
		UID:  UID1,
		Note: note,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respCreate.NID)

	respCreate, err = st.NoteClient.CreateNote(ctx, &notes.CreateNoteRequest{
		UID:  UID1,
		Note: note,
	})

	require.Error(t, err)
	assert.Empty(t, respCreate)
	assert.ErrorContains(t, err, "note with that title already exists")
}

func TestCreate_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name          string
		note          *notes.Note
		expectedError string
	}{
		{
			name: "Empty content",
			note: &notes.Note{
				Title:    gofakeit.Name(),
				Content:  "",
				Duration: durationpb.New(time.Minute * 10),
			},
			expectedError: "empty content field",
		},
		{
			name: "Empty title",
			note: &notes.Note{
				Title:    "",
				Content:  gofakeit.HackerPhrase(),
				Duration: durationpb.New(time.Minute * 10),
			},
			expectedError: "empty title field",
		},
		{
			name: "Both empty, but err content",
			note: &notes.Note{
				Title:    "",
				Content:  "",
				Duration: durationpb.New(time.Minute * 10),
			},
			expectedError: "empty content field",
		},
		{
			name: "Ivalid duration value",
			note: &notes.Note{
				Title:    gofakeit.Name(),
				Content:  gofakeit.HackerPhrase(),
				Duration: durationpb.New(time.Minute * 9),
			},
			expectedError: "invalid time duration",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.NoteClient.CreateNote(ctx, &notes.CreateNoteRequest{
				UID:  UID2,
				Note: tt.note,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

func TestCRUD_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	note := &notes.Note{
		Title:    gofakeit.Name(),
		Content:  gofakeit.HackerPhrase(),
		Duration: durationpb.New(time.Minute * 10),
	}

	respCreate, err := st.NoteClient.CreateNote(ctx, &notes.CreateNoteRequest{
		UID:  UID1,
		Note: note,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respCreate.NID)

	respGet, err := st.NoteClient.GetNote(ctx, &notes.NoteIDRequest{
		UID: UID1,
		NID: respCreate.NID,
	})

	assert.NoError(t, err)
	require.Equal(t, note.Title, respGet.Title)
	require.Equal(t, note.Content, respGet.Content)

	newNote := &notes.Note{
		Title:    gofakeit.Name(),
		Content:  gofakeit.HackerPhrase(),
		Duration: durationpb.New(time.Minute * 10),
	}

	respUpdate, err := st.NoteClient.UpdateNote(ctx, &notes.UpdateNoteRequest{
		UID:  UID1,
		NID:  respCreate.NID,
		Note: newNote,
	})

	assert.NoError(t, err)
	require.Equal(t, respCreate.NID, respUpdate.NID)

	respNewGet, err := st.NoteClient.GetNote(ctx, &notes.NoteIDRequest{
		NID: respUpdate.NID,
		UID: UID1,
	})

	assert.NoError(t, err)
	require.Equal(t, respNewGet.Title, newNote.Title)
	require.Equal(t, respNewGet.Content, newNote.Content)
	require.Equal(t, respNewGet.Duration.AsDuration(), newNote.Duration.AsDuration())
}

func TestUpdate_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	note := &notes.Note{
		Title:    gofakeit.Name(),
		Content:  gofakeit.HackerPhrase(),
		Duration: durationpb.New(time.Minute * 10),
	}

	respCreate, err := st.NoteClient.CreateNote(ctx, &notes.CreateNoteRequest{
		UID:  UID2,
		Note: note,
	})

	assert.NoError(t, err)
	require.NotEmpty(t, respCreate.NID)

	tests := []struct {
		name          string
		note          *notes.Note
		expectedError string
	}{
		{
			name: "Full empty update",
			note: &notes.Note{
				Title:    "",
				Content:  "",
				Duration: durationpb.New(0),
			},
			expectedError: "invalid update request",
		},
		{
			name: "Invalid time duration value",
			note: &notes.Note{
				Title:    "",
				Content:  "some Content",
				Duration: durationpb.New(time.Minute * 3),
			},
			expectedError: "invalid time duration",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.NoteClient.UpdateNote(ctx,
				&notes.UpdateNoteRequest{
					UID:  UID2,
					NID:  respCreate.NID,
					Note: tt.note,
				},
			)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

func TestDelete_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	note := &notes.Note{
		Title:    gofakeit.Name(),
		Content:  gofakeit.HackerPhrase(),
		Duration: durationpb.New(time.Minute * 10),
	}

	respCreate, err := st.NoteClient.CreateNote(ctx, &notes.CreateNoteRequest{
		UID:  UID2,
		Note: note,
	})

	assert.NoError(t, err)
	require.NotEmpty(t, respCreate.NID)

	respGet, err := st.NoteClient.GetNote(ctx, &notes.NoteIDRequest{
		UID: UID2,
		NID: respCreate.NID,
	})

	assert.NoError(t, err)
	require.Equal(t, note.Title, respGet.Title)
	require.Equal(t, note.Content, respGet.Content)
	require.Equal(t, note.Duration.AsDuration(), respGet.Duration.AsDuration())

	respDelete, err := st.NoteClient.DeleteNote(ctx, &notes.NoteIDRequest{
		UID: UID2,
		NID: respCreate.NID,
	})

	assert.NoError(t, err)
	require.Equal(t, respDelete.NID, respCreate.NID)

	_, err = st.NoteClient.GetNote(ctx, &notes.NoteIDRequest{
		UID: UID2,
		NID: respCreate.NID,
	})

	assert.Error(t, err)
	require.Contains(t, err.Error(), "note not found")
}

func TestListNotes_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	notesCount := 4
	notesMap := make(map[int64]*notes.Note, notesCount)
	notesID := make(map[int64]interface{}, notesCount)

	for range notesCount {
		note := &notes.Note{
			Title:    gofakeit.Name(),
			Content:  gofakeit.HackerPhrase(),
			Duration: durationpb.New(time.Minute * 10),
		}

		respCreate, err := st.NoteClient.CreateNote(ctx, &notes.CreateNoteRequest{
			UID:  UID2,
			Note: note,
		})
		assert.NoError(t, err)

		notesMap[respCreate.NID] = note
		notesID[respCreate.NID] = true
	}

	respListID, err := st.NoteClient.ListUserNotesID(ctx, &notes.UserIDRequest{
		UID: UID2,
	})

	assert.NoError(t, err)
	require.LessOrEqual(t, len(notesID), len(respListID.NoteIDs))

	for _, id := range respListID.NoteIDs {
		delete(notesID, id)
	}

	require.Equal(t, 0, len(notesID))

	respListNotes, err := st.NoteClient.ListUserNotes(ctx, respListID)

	assert.NoError(t, err)
	require.LessOrEqual(t, len(respListID.NoteIDs), len(respListNotes.Notes))

	for _, note := range respListNotes.Notes {
		if expectedNote, ok := notesMap[note.NID]; ok {
			require.Equal(t, expectedNote.Title, note.Note.Title)
			require.Equal(t, expectedNote.Content, note.Note.Content)
			require.Equal(t, expectedNote.Duration.AsDuration(), note.Note.Duration.AsDuration())
		}
	}

}
