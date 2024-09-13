package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *SUser) GetAllPermissions() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		department := common.ParseQueryCriteriaToIntArray(c.Query("department"))
		query := c.Query("query")

		proto, err := s.RoleClient.GetAllPermissions(ctx, &userProto.GetAllPermissionsRequest{
			Department: department,
			Query:      query,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildAllPermissionsResponse(proto.GetResult()))
	}
}
