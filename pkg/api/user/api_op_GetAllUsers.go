package user

import (
	"context"
	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"net/http"
)

func (s *SUser) GetAllUsers() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		perPage := common.ConvertStringToInt(c.Query("perPage"))
		currentPage := common.ConvertStringToInt(c.Query("currentPage"))
		department := common.ParseQueryCriteriaToIntArray(c.Query("department"))
		role := common.ParseQueryCriteriaToIntArray(c.Query("role"))
		query := c.Query("query")
		fromDate := c.Query("fromDate")
		toDate := c.Query("toDate")

		proto, err := s.UserClient.GetAllUsers(ctx, &userProto.GetAllUsersRequest{
			Params: &userProto.UPaginationParams{
				PerPage:     perPage,
				CurrentPage: currentPage,
			},
			Role:       role,
			Department: department,
			Query:      query,
			FromDate:   fromDate,
			ToDate:     toDate,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		response := openapi.GetAllUsersResponse{
			Pagination: buildPaginationResponse(proto.GetPagination()),
		}

		if proto.GetResult() == nil {
			response.Result = []openapi.User{}
		} else {
			response.Result = buildAllUsersResponse(proto.GetResult())
		}

		c.JSON(http.StatusOK, response)
	}
}
