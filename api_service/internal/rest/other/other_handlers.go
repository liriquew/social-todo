package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	friends_grpc "github.com/liriquew/social-todo/api_service/internal/clients/friendsgrpc"
	notes_grpc "github.com/liriquew/social-todo/api_service/internal/clients/notesgrpc"
)

type GeneralAPI interface {
	ListLastNotes(*gin.Context)
}

type General struct {
	log           *slog.Logger
	notesClient   *notes_grpc.Client
	friendsClient *friends_grpc.Client
}

func New(log *slog.Logger, notesClient *notes_grpc.Client, friendsClient *friends_grpc.Client) *General {
	return &General{
		log:           log,
		notesClient:   notesClient,
		friendsClient: friendsClient,
	}
}

func (a *General) ListLastNotes(c *gin.Context) {
	uid := c.Value("uid").(int64)
	if uid <= 0 {
		c.Status(http.StatusUnauthorized)
		return
	}

	offset, err := strconv.ParseInt(c.Query("offset"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}

	if offset < 0 {
		c.Status(http.StatusBadRequest)
		return
	}
	if limit > 10 {
		c.String(http.StatusBadRequest, "limit value is too big")
	}

	FIDs, err := a.friendsClient.ListFriends(c, uid)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(FIDs) == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	notes, err := a.notesClient.ListUsersNotes(c, FIDs, offset, limit)
	if err != nil {
		if errors.Is(err, notes_grpc.ErrNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, notes)
}
