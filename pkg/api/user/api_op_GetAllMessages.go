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
)

func (s *SUser) GetAllMessages() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		perPage := common.ConvertStringToInt(c.Query("perPage"))
		currentPage := common.ConvertStringToInt(c.Query("currentPage"))

		channelName := auth.GetAccountMetadataFromToken(
			strings.TrimPrefix(c.GetHeader(meta.AuthorizationHeaderKey), "Bearer "),
		).Email
		channelName = strings.ReplaceAll(channelName, ".", "-")
		channelName = "messages:" + strings.ReplaceAll(channelName, "@", "-")

		startIndex := int64((currentPage - 1) * perPage)
		stopIndex := startIndex + int64(perPage) - 1

		result := []map[string]interface{}{}
		for _, msg := range common.RdbInstance.ZRange(ctx, channelName, startIndex, stopIndex).Val() {
			var jsonMessage map[string]interface{}
			json.Unmarshal([]byte(msg), &jsonMessage)
			result = append(result, jsonMessage)
		}

		totalItems := int32(common.RdbInstance.ZCard(ctx, channelName).Val())
		lastPage := (totalItems + int32(perPage) - 1) / int32(perPage)

		c.JSON(http.StatusOK, map[string]interface{}{
			"result": result,
			"pagination": openapi.Pagination{
				Total:       totalItems,
				PerPage:     perPage,
				CurrentPage: currentPage,
				LastPage:    lastPage,
				From:        int32(startIndex + 1),
				To:          int32(min(stopIndex+1, int64(totalItems))),
			},
		})
	}
}
