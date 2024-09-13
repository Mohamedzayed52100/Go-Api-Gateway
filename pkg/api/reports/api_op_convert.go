package reports

import (
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	reportsProto "github.com/goplaceapp/shared-protobufs/reports/go_out"
)

func buildCustomizeReportResponse(proto []*reportsProto.CustomizeReportResponse) []openapi.CustomizeReport {
	var result []openapi.CustomizeReport

	for _, item := range proto {
		var reservationTags []openapi.ReportTag
		for _, tag := range item.GetReservationTags() {
			reservationTags = append(reservationTags, openapi.ReportTag{
				Name:  tag.GetName(),
				Color: tag.GetColor(),
			})
		}

		result = append(result, openapi.CustomizeReport{
			Id:              item.GetId(),
			Name:            item.GetName(),
			Date:            item.GetDate(),
			Time:            item.GetTime(),
			Branch:          item.GetBranch(),
			ReservationTags: reservationTags,
			Spent:           item.GetSpent(),
			SeatingArea:     item.GetSeatingArea(),
			Shift:           item.GetShift(),
			BookedVia:       item.GetBookedVia(),
			Status: openapi.ReportTag{
				Name:  item.GetStatus().GetName(),
				Color: item.GetStatus().GetColor(),
			},
			GuestsNumber: item.GetGuestsNumber(),
		})
	}

	return result
}

func buildGuestOverviewReportResponse(proto []*reportsProto.OverviewGuestReport) []openapi.OverviewGuestReport {
	var result []openapi.OverviewGuestReport

	for _, item := range proto {
		result = append(result, openapi.OverviewGuestReport{
			Period:          item.GetPeriod(),
			Repeater:        item.GetRepeater(),
			Vip:             item.GetVip(),
			SpecialOccasion: item.GetSpecialOccasion(),
		})
	}

	return result
}

func BuildPaginationResponse(proto *reportsProto.Pagination) openapi.Pagination {
	return openapi.Pagination{
		Total:       proto.GetTotal(),
		PerPage:     proto.GetPerPage(),
		CurrentPage: proto.GetCurrentPage(),
		LastPage:    proto.GetLastPage(),
		From:        proto.GetFrom(),
		To:          proto.GetTo(),
	}
}

func buildGuestReportResponse(proto *reportsProto.GetGuestReportResponse) openapi.GetGuestReportResponse {
	event := proto.GetSpecialOccasion().GetEvents()

	var events []openapi.EventCount
	for _, item := range event {
		events = append(events, openapi.EventCount{
			Count: item.GetCount(),
			Type:  item.GetType(),
		})
	}

	return openapi.GetGuestReportResponse{
		Vip:      proto.GetVip(),
		Repeater: proto.GetRepeater(),
		SpecialOccasion: openapi.GuestReportSpecialOccasionCount{
			Count:      proto.GetSpecialOccasion().GetCount(),
			TotalCount: proto.GetSpecialOccasion().GetTotalCount(),
			Events:     events,
		},
		ReservedBy: openapi.GuestReportReservedByCount{
			Whatsapp: proto.GetReservedBy().GetWhatsapp(),
			Goplace:  proto.GetReservedBy().GetGoplace(),
			Other:    proto.GetReservedBy().GetOther(),
		},
	}
}

func buildTopSpendersFinancialReportResponse(proto *reportsProto.GetFinancialTopSpendersReportResponse) []openapi.TopSpendersTransaction {
	var result []openapi.TopSpendersTransaction

	for _, item := range proto.GetResult() {
		result = append(result, openapi.TopSpendersTransaction{
			Name:    item.GetName(),
			GuestId: item.GetGuestId(),
			Amount:  item.GetAmount(),
		})
	}

	return result
}

func buildOverviewReservationsReport(proto []*reportsProto.ReservationData) []openapi.OverviewReservationsReport {
	var result []openapi.OverviewReservationsReport

	for _, item := range proto {
		result = append(result, openapi.OverviewReservationsReport{
			Time:             item.GetTime(),
			TotalReservation: item.GetTotalReservations(),
			TotalNoShow:      item.GetTotalNoShows(),
		})
	}

	return result
}
