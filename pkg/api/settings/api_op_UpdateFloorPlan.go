package settings

import (
	"context"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) UpdateFloorPlan() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.UpdateFloorPlanRequest
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

		id := c.Param("id")
		proto, err := s.FloorPlanClient.UpdateFloorPlan(ctx, &settingsProto.UpdateFloorPlanRequest{
			Params: &settingsProto.FloorPlanParams{
				Id:            common.ConvertStringToInt(id),
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

		c.JSON(200, buildFloorPlanResponse(proto.GetResult()))
	}
}
