package settings

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) GetAllReservationStatuses() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var (
			proto     *settingsProto.GetAllReservationStatusesResponse
			err       error
			shortData = c.Query("shortData")
			branchIds = common.ParseQueryCriteriaToIntArray(c.Query("branchIds"))
		)

		if shortData == "true" {
			proto, err = s.ReservationStatusClient.GetAllReservationStatusesShort(ctx, &settingsProto.GetAllReservationStatusesRequest{
				BranchIds: branchIds,
			})
			if err != nil {
				common.HandleGrpcError(c, err)
				return
			}
		} else {
			proto, err = s.ReservationStatusClient.GetAllReservationStatuses(ctx, &settingsProto.GetAllReservationStatusesRequest{
				BranchIds: branchIds,
			})
			if err != nil {
				common.HandleGrpcError(c, err)
				return
			}
		}

		c.JSON(http.StatusAccepted, buildAllBranchReservationStatusesResponse(proto.GetResult()))
	}
}
