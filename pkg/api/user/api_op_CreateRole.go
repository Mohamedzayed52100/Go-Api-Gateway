package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *SUser) CreateRole() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.CreateRoleRequest
		common.ResolveRequestBody(c, &input)

		proto, err := s.RoleClient.CreateRole(ctx, &userProto.CreateRoleRequest{
			Params: &userProto.RoleParams{
				DisplayName: input.DisplayName,
				Department:  input.Department,
				Permissions: *input.Permissions,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildRoleResponse(proto.GetResult()))
	}
}
