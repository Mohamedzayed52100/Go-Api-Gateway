package settings

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) CreateFloorPlan() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.CreateFloorPlanRequest
		common.ResolveRequestBody(c, &input)

		elements := []*settingsProto.FloorPlanElementParams{}
		for _, i := range input.Elements {
			elements = append(elements, &settingsProto.FloorPlanElementParams{
				ElementType: i.ElementType,
				TableId:     i.TableId,
				Property: &settingsProto.FloorPlanElementProperty{
					Shape:  i.Property.Shape,
					Width:  i.Property.Width,
					Height: i.Property.Height,
					Angle:  i.Property.Angle,
					X:      i.Property.X,
					Y:      i.Property.Y,
				},
			})
		}

		proto, err := s.FloorPlanClient.CreateFloorPlan(ctx, &settingsProto.CreateFloorPlanRequest{
			Params: &settingsProto.FloorPlanParams{
				Name:          input.Name,
				SeatingAreaId: input.SeatingAreaId,
				Scale:         input.Scale,
				Elements:      elements,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildFloorPlanResponse(proto.GetResult()))
	}
}
