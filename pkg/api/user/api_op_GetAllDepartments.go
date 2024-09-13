package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *SUser) GetAllDepartments() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		perPage := common.ConvertStringToInt(c.Query("perPage"))
		currentPage := common.ConvertStringToInt(c.Query("currentPage"))

		proto, err := s.DepartmentClient.GetAllDepartments(ctx, &userProto.GetAllDepartmentsRequest{
			Params: &userProto.UPaginationParams{
				PerPage:     perPage,
				CurrentPage: currentPage,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		res := openapi.GetAllDepartmentsResponse{
			Result: buildAllDepartmentsResponse(proto.GetResult()),
		}

		if perPage == 0 && currentPage == 0 {
			res.Pagination = nil
		} else {
			pagination := buildPaginationResponse(proto.GetPagination())
			res.Pagination = &pagination
		}
		c.JSON(http.StatusOK, res)
	}
}
