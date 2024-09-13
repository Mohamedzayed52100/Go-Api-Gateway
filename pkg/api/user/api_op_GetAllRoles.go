package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *SUser) GetAllRoles() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		perPage := common.ConvertStringToInt(c.Query("perPage"))
		currentPage := common.ConvertStringToInt(c.Query("currentPage"))
		department := common.ParseQueryCriteriaToIntArray(c.Query("department"))
		query := c.Query("query")

		proto, err := s.RoleClient.GetAllRoles(ctx, &userProto.GetAllRolesRequest{
			Params: &userProto.UPaginationParams{
				PerPage:     perPage,
				CurrentPage: currentPage,
			},
			Department: department,
			Query:      query,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		res := openapi.GetAllRolesResponse{
			Pagination: buildPaginationResponse(proto.GetPagination()),
			Result:     buildAllRolesResponse(proto.GetResult()),
		}

		if proto.GetResult() == nil {
			res.Result = []openapi.Role{}
		}

		c.JSON(http.StatusOK, res)
	}
}
