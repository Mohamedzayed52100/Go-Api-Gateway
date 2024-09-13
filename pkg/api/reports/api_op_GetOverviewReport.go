package reports

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	reports "github.com/goplaceapp/shared-protobufs/reports/go_out"
)

func (s *SReports) GetOverviewReport() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		fromDate := c.Query("fromDate")
		toDate := c.Query("toDate")
		profileTagIds := c.Query("profileTagIds")
		reservationTagIds := c.Query("reservationTagIds")
		bookedVia := c.Query("bookedVia")
		period := c.Query("period")
		branchIds := c.Query("branchIds")

		proto, err := s.OverviewReportClient.GetOverviewReport(ctx, &reports.ReportRequestParam{
			FromDate:          fromDate,
			ToDate:            toDate,
			ProfileTagIds:     common.ParseQueryCriteriaToIntArray(profileTagIds),
			ReservationTagIds: common.ParseQueryCriteriaToIntArray(reservationTagIds),
			BookedVia:         common.ParseQueryCriteriaToStringArray(bookedVia),
			Period:            period,
			BranchIds:         common.ParseQueryCriteriaToIntArray(branchIds),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		var events []openapi.EventCount
		for _, item := range proto.GetSpecialOccasion().GetEvents() {
			events = append(events, openapi.EventCount{
				Count: item.GetCount(),
				Type:  item.GetType(),
			})
		}

		c.JSON(http.StatusOK, &openapi.GetOverviewReportResponse{
			Reservations: openapi.OverviewReport{
				Count:            proto.GetReservations().GetCount(),
				TotalGuestsCount: proto.GetReservations().GetTotalGuestsCount(),
				Percentage:       proto.GetReservations().GetPercentage(),
			},
			NoShow: openapi.OverviewReport{
				Count:            proto.GetNoShow().GetCount(),
				Percentage:       proto.GetNoShow().GetPercentage(),
				TotalGuestsCount: proto.GetNoShow().GetTotalGuestsCount(),
			},
			TopTimes: proto.GetTopTimes(),
			Vip: openapi.OverviewGuestTag{
				Count:      proto.GetVip().GetCount(),
				TotalCount: proto.GetVip().GetTotalCount(),
			},
			Repeater: openapi.OverviewGuestTag{
				Count:      proto.GetRepeater().GetCount(),
				TotalCount: proto.GetRepeater().GetTotalCount(),
			},
			SpecialOccasion: openapi.GuestReportSpecialOccasionCount{
				Count:      proto.GetSpecialOccasion().GetCount(),
				TotalCount: proto.GetSpecialOccasion().GetTotalCount(),
				Events:     events,
			},
		})
	}
}
