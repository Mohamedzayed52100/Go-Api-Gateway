package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *SUser) UpdateRole() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var (
			input openapi.UpdateRoleRequest
			proto *userProto.UpdateRoleResponse
			err   error
		)
		common.ResolveRequestBody(c, &input)
		id := common.ConvertStringToInt(c.Param("roleId"))
		if input.Permissions != nil {
			proto, err = s.RoleClient.UpdateRole(ctx, &userProto.UpdateRoleRequest{
				Params: &userProto.RoleParams{
					Id:               id,
					Name:             input.Name,
					DisplayName:      input.DisplayName,
					Department:       input.Department,
					EmptyPermissions: false,
					Permissions:      *input.Permissions,
				},
			})
		} else {
			proto, err = s.RoleClient.UpdateRole(ctx, &userProto.UpdateRoleRequest{
				Params: &userProto.RoleParams{
					Id:               id,
					Name:             input.Name,
					DisplayName:      input.DisplayName,
					Department:       input.Department,
					EmptyPermissions: true,
				},
			})
		}
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildRoleResponse(proto.GetResult()))
	}
}
