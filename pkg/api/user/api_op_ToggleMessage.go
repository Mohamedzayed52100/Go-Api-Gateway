package user

import (
	"context"
	"encoding/json"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/goplaceapp/goplace-common/pkg/auth"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	"github.com/redis/go-redis/v9"
)

func (s *SUser) ToggleMessage() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		id := c.Param("messageId")
		rdb := common.RdbInstance

		channelName := auth.GetAccountMetadataFromToken(ctx.Value(meta.AuthorizationContextKey.String()).(string)).Email
		channelName = strings.ReplaceAll(channelName, ".", "-")
		channelName = "messages:" + strings.ReplaceAll(channelName, "@", "-")

		allValues := rdb.ZRange(ctx, channelName, 0, -1).Val()

		var oldValue string
		for _, v := range allValues {
			if strings.Contains(v, id) {
				oldValue = v
				break
			}
		}

		var jsonMessage map[string]interface{}
		if oldValue != "" {
			err := json.Unmarshal([]byte(oldValue), &jsonMessage)
			if err != nil {
				c.JSON(http.StatusBadRequest, openapi.ToggleMessageResponse{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, openapi.ToggleMessageResponse{
				Code:    http.StatusBadRequest,
				Message: "Message not found",
			})
			return
		}
		jsonMessage["seen"] = !jsonMessage["seen"].(bool)

		if err := rdb.ZRem(ctx, channelName, oldValue).Err(); err != nil {
			c.JSON(http.StatusBadRequest, openapi.ToggleMessageResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
			return
		}

		newValue, err := json.Marshal(jsonMessage)
		if err != nil {
			c.JSON(http.StatusBadRequest, openapi.ToggleMessageResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
			return
		}

		if err := rdb.ZAdd(ctx, channelName, redis.Z{
			Member: string(newValue),
		}).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, openapi.ToggleMessageResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}

		if jsonMessage["seen"].(bool) {
			c.JSON(http.StatusOK, openapi.ToggleMessageResponse{
				Code:    http.StatusOK,
				Message: "Message marked as read",
			})
		} else {
			c.JSON(http.StatusOK, openapi.ToggleMessageResponse{
				Code:    http.StatusOK,
				Message: "Message marked as unread",
			})
		}
	}
}
