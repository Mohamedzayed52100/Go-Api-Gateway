package settings

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) CreateShift() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		editMode := c.Query("editMode")
		deleteMode := c.Query("deleteMode")
		date := c.Query("date")

		var input openapi.CreateShiftRequest
		common.ResolveRequestBody(c, &input)

		if input.EndDate == "infinity" {
			input.EndDate = meta.InfiniteDate
		}

		turnover := make([]*settingsProto.TurnoverParams, 0)
		for _, t := range input.Turnover {
			turnover = append(turnover, &settingsProto.TurnoverParams{
				ShiftId:      t.ShiftId,
				GuestsNumber: t.GuestsNumber,
				TurnoverTime: t.TurnoverTime,
			})
		}

		pacing := make([]*settingsProto.PacingParams, 0)
		for _, t := range input.Pacing {
			pacing = append(pacing, &settingsProto.PacingParams{
				Hour:          t.Hour,
				Capacity:      t.Capacity,
				SeatingAreaId: t.SeatingAreaId,
			})
		}

		daysToRepeat := ""
		for _, day := range input.DaysToRepeat {
			if day == input.DaysToRepeat[len(input.DaysToRepeat)-1] {
				daysToRepeat += day
			} else {
				daysToRepeat += day + ","
			}
		}

		proto, err := s.ShiftClient.CreateShift(ctx, &settingsProto.CreateShiftRequest{
			Params: &settingsProto.ShiftParams{
				BranchId:       input.BranchId,
				Name:           input.Name,
				From:           input.From,
				To:             input.To,
				StartDate:      input.StartDate,
				EndDate:        input.EndDate,
				TimeInterval:   input.TimeInterval,
				FloorPlanId:    input.FloorPlanId,
				SeatingAreaIds: input.SeatingAreaIds,
				CategoryId:     input.CategoryId,
				MinGuests:      input.MinGuests,
				MaxGuests:      input.MaxGuests,
				DaysToRepeat:   daysToRepeat,
				Turnover:       turnover,
				Pacing:         pacing,
				EditMode:       editMode,
				DeleteMode:     deleteMode,
				Date:           date,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildShiftResponse(proto.GetResult()))
	}
}
