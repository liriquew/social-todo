package friends

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	friends_grpc "github.com/liriquew/social-todo/api_service/internal/clients/friendsgrpc"
	"github.com/liriquew/social-todo/api_service/internal/models"
	"github.com/liriquew/social-todo/api_service/pkg/logger/sl"
)

type FriendsAPI interface {
	AddFriend(*gin.Context)
	RemoveFriend(*gin.Context)
	ListFriend(*gin.Context)
}

type Friends struct {
	log           *slog.Logger
	friendsClient *friends_grpc.Client
}

func New(log *slog.Logger, friendsClient *friends_grpc.Client) *Friends {
	return &Friends{
		log:           log,
		friendsClient: friendsClient,
	}
}

func (f *Friends) AddFriend(c *gin.Context) {
	uid := c.Value("uid").(int64)
	if uid <= 0 {
		c.Status(http.StatusUnauthorized)
		return
	}

	var FID models.FriendID
	if err := c.ShouldBindJSON(&FID); err != nil {
		f.log.Warn("bad json", sl.Err(err))
		c.String(http.StatusBadRequest, "bad json idk")
		return
	}

	err := f.friendsClient.AddFriend(c, uid, FID.FID)
	if err != nil {
		f.log.Warn("bad json", sl.Err(err))
		c.String(http.StatusBadRequest, "bad json idk")
		return
	}

	c.Status(http.StatusOK)
}

func (f *Friends) RemoveFriend(c *gin.Context) {
	uid := c.Value("uid").(int64)
	if uid <= 0 {
		c.Status(http.StatusUnauthorized)
		return
	}

	var FID models.FriendID
	if err := c.ShouldBindJSON(&FID); err != nil {
		f.log.Warn("bad json", sl.Err(err))
		c.String(http.StatusBadRequest, "bad json idk")
		return
	}

	err := f.friendsClient.RemoveFriend(c, uid, FID.FID)
	if err != nil {
		f.log.Warn("bad json", sl.Err(err))
		c.String(http.StatusBadRequest, "bad json idk")
		return
	}

	c.Status(http.StatusOK)
}

func (f *Friends) ListFriend(c *gin.Context) {
	uid := c.Value("uid").(int64)
	if uid <= 0 {
		c.Status(http.StatusUnauthorized)
		return
	}

	FIDs, err := f.friendsClient.ListFriends(c, uid)
	if err != nil {
		f.log.Warn("bad json", sl.Err(err))
		c.String(http.StatusBadRequest, "bad json idk")
		return
	}

	c.JSON(http.StatusOK, FIDs)
}
