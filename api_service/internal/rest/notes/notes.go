package notes

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	notes_grpc "github.com/liriquew/social-todo/api_service/internal/clients/notesgrpc"
	"github.com/liriquew/social-todo/api_service/internal/models"
	"github.com/liriquew/social-todo/api_service/pkg/logger/sl"
)

type NotesAPI interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// this is api
type Notes struct {
	log         *slog.Logger
	notesClient *notes_grpc.Client
}

func New(log *slog.Logger, notesClient *notes_grpc.Client) *Notes {
	return &Notes{
		log:         log,
		notesClient: notesClient,
	}
}

// TODO: implement endpoints
func (n *Notes) Create(c *gin.Context) {
	uid := c.Value("uid").(int64)
	if uid <= 0 {
		c.Status(http.StatusUnauthorized)
		return
	}

	var note *models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		n.log.Warn("bad json", sl.Err(err))
		c.String(http.StatusBadRequest, "bad json idk")
		return
	}

	NoteID, err := n.notesClient.Create(c, uid, note)
	if err != nil {
		n.log.Warn("error:", sl.Err(err))

		if errors.Is(err, notes_grpc.ErrAlreadyExists) {
			c.Status(http.StatusConflict)
			return
		}
		if errors.Is(err, notes_grpc.ErrInvalidArgument) {
			c.Status(http.StatusBadRequest)
			return
		}

		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"note_id": NoteID,
	})
}

func (n *Notes) Get(c *gin.Context) {
	uid := c.Value("uid").(int64)
	if uid <= 0 {
		c.Status(http.StatusUnauthorized)
		return
	}

	noteId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if noteId <= 0 || err != nil {
		n.log.Warn("ERROR GET parseInt", sl.Err(err))
		c.Status(http.StatusBadRequest)
	}

	note, err := n.notesClient.Get(c, uid, noteId)
	if err != nil {
		n.log.Warn("error:", sl.Err(err))

		// TODO: check err not exists
		if errors.Is(err, notes_grpc.ErrNotFound) {
			c.Status(http.StatusNotFound)
			return
		}

		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, note)
}

func (n *Notes) Update(c *gin.Context) {
	uid := c.Value("uid").(int64)
	if uid <= 0 {
		c.Status(http.StatusUnauthorized)
		return
	}

	noteId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if noteId <= 0 || err != nil {
		n.log.Warn("ERROR GET parseInt", sl.Err(err))
		c.Status(http.StatusOK)
	}

	var note *models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		n.log.Warn("bad json", sl.Err(err))
		c.String(http.StatusBadRequest, "bad json idk")
		return
	}

	err = n.notesClient.Update(c, uid, noteId, note)
	if err != nil {
		n.log.Warn("error:", sl.Err(err))

		if errors.Is(err, notes_grpc.ErrInvalidArgument) {
			c.Status(http.StatusBadRequest)
			return
		}
		if errors.Is(err, notes_grpc.ErrAlreadyExists) {
			c.Status(http.StatusConflict)
			return
		}
		if errors.Is(err, notes_grpc.ErrNotFound) {
			c.Status(http.StatusNotFound)
			return
		}

		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (n *Notes) Delete(c *gin.Context) {
	uid := c.Value("uid").(int64)
	if uid <= 0 {
		c.Status(http.StatusUnauthorized)
		return
	}

	noteId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if noteId <= 0 || err != nil {
		n.log.Warn("ERROR GET parseInt", sl.Err(err))
		c.Status(http.StatusOK)
	}

	err = n.notesClient.Delete(c, uid, noteId)
	if err != nil {
		n.log.Warn("error:", sl.Err(err))

		c.Status(http.StatusInternalServerError)
	}

	c.Status(http.StatusOK)
}
