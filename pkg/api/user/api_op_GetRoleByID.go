package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *SUser) GetRoleByID() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		id := common.ConvertStringToInt(c.Param("roleId"))

		proto, err := s.RoleClient.GetRoleByID(ctx, &userProto.GetRoleByIDRequest{
			Id: id,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildRoleResponse(proto.GetResult()))
	}
}
