package customerservice

import (
	"time"

	cssProto "github.com/goplaceapp/goplace-customer-service/api/v1"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
)

func buildEmployeesReservationsCountResponse(reservationsCount []*cssProto.EmployeeReservationCount) []openapi.EmployeeReservationCount {
	result := []openapi.EmployeeReservationCount{}
	for _, reservation := range reservationsCount {
		result = append(result, openapi.EmployeeReservationCount{
			EmployeeId: reservation.GetEmployeeId(),
			Name:       reservation.GetName(),
			Avatar:     reservation.GetAvatar(),
			Branches:   buildBranchReservationsResponse(reservation.GetBranches()),
		})
	}

	return result
}

func buildBranchReservationsResponse(branches []*cssProto.BranchReservations) []openapi.BranchReservations {
	response := []openapi.BranchReservations{}
	for _, branch := range branches {
		response = append(response, openapi.BranchReservations{
			BranchId:   branch.GetBranchId(),
			BranchName: branch.GetBranchName(),
			Count:      branch.GetCount(),
		})
	}

	return response
}

// Guest
func buildGuestResponse(guest *cssProto.CGuest) openapi.CGuest {
	return openapi.CGuest{
		Id:          guest.GetId(),
		FirstName:   guest.GetFirstName(),
		LastName:    guest.GetLastName(),
		PhoneNumber: guest.GetPhoneNumber(),
	}
}

// Branch
func buildBranchResponse(branch *cssProto.CBranch) openapi.CBranch {
	return openapi.CBranch{
		Id:   int32(branch.GetId()),
		Name: branch.GetName(),
	}
}

// Inquiry Status
func buildInquiryStatusResponse(status *cssProto.InquiryStatus) openapi.InquiryStatus {
	return openapi.InquiryStatus{
		Id:        int32(status.GetId()),
		Name:      status.GetName(),
		Color:     status.GetColor(),
		CreatedAt: status.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt: status.GetUpdatedAt().AsTime().Format(time.RFC822),
	}
}

// Inquiry Solution
func buildInquirySolutionResponse(solution *cssProto.InquirySolution) openapi.InquirySolution {
	return openapi.InquirySolution{
		Id:        int32(solution.GetId()),
		Creator:   buildCreatorResponse(solution.Creator),
		Solution:  solution.GetSolution(),
		CreatedAt: solution.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt: solution.GetUpdatedAt().AsTime().Format(time.RFC822),
	}
}

// Inquiry Type
func buildInquiryTypeResponse(inquiryType *cssProto.InquiryType) openapi.InquiryType {
	return openapi.InquiryType{
		Id:        int32(inquiryType.GetId()),
		Name:      inquiryType.GetName(),
		CreatedAt: inquiryType.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt: inquiryType.GetUpdatedAt().AsTime().Format(time.RFC822),
	}
}

// Inquiry Creator
func buildCreatorResponse(creator *cssProto.CCreator) openapi.Creator {
	return openapi.Creator{
		Id:          creator.GetId(),
		FirstName:   creator.GetFirstName(),
		LastName:    creator.GetLastName(),
		PhoneNumber: creator.GetPhoneNumber(),
		Email:       creator.GetEmail(),
		Avatar:      creator.GetAvatar(),
		Role:        creator.GetRole(),
	}
}

// Inquiry InquiryComment
func buildInquiryCommentResponse(comment *cssProto.InquiryComment) openapi.InquiryComment {
	res := openapi.InquiryComment{
		Id:        int32(comment.GetId()),
		Creator:   buildCreatorResponse(comment.GetCreator()),
		Comment:   comment.GetComment(),
		CreatedAt: comment.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt: comment.GetUpdatedAt().AsTime().Format(time.RFC822),
	}

	if comment.Inquiry != nil {
		inquiry := buildInquiryResponse(comment.GetInquiry())
		res.Inquiry = &inquiry
	} else {
		res.Inquiry = nil
	}

	return res
}

// All Inquiry comments
func buildAllInquiryCommentsResponse(comments []*cssProto.InquiryComment) []openapi.InquiryComment {
	response := []openapi.InquiryComment{}
	for _, comment := range comments {
		response = append(response, buildInquiryCommentResponse(comment))
	}

	return response
}

func buildTopEmployeesResponse(employees []*cssProto.TopEmployee) []openapi.TopEmployee {
	response := []openapi.TopEmployee{}
	for _, employee := range employees {
		response = append(response, buildEmployeeResponse(employee))
	}

	return response
}

func buildEmployeeResponse(employee *cssProto.TopEmployee) openapi.TopEmployee {
	return openapi.TopEmployee{
		Id:                employee.GetId(),
		Name:              employee.GetName(),
		Avatar:            employee.GetAvatar(),
		Role:              employee.GetRole(),
		ReservationsCount: employee.GetReservationsCount(),
		InquiriesCount:    employee.GetInquiriesCount(),
		AvgTime:           employee.GetAvgTime(),
	}
}

// buildShiftReponse
func buildShiftResponse(shift *cssProto.CShift) openapi.CShift {
	return openapi.CShift{
		Id:   int32(shift.GetId()),
		Name: shift.GetName(),
		Cast: []openapi.CCast{},
	}
}

// buildReservationResponse
func buildReservationResponse(reservation *cssProto.CReservation) openapi.CReservation {
	res := openapi.CReservation{
		Id:           int32(reservation.GetId()),
		Date:         reservation.GetDate().AsTime().Format(time.DateOnly),
		Branch:       buildBranchResponse(reservation.GetBranch()),
		Shift:        buildShiftResponse(reservation.GetShift()),
		ReservedVia:  reservation.GetReservedVia(),
		GuestsNumber: int32(reservation.GetGuestsNumber()),
		Tables:       buildInquiryReservationTablesProto(reservation.GetTables()),
	}

	if reservation.GetCheckIn() != nil {
		res.CheckIn = reservation.GetCheckIn().AsTime().Format(time.RFC822)
	} else {
		res.CheckIn = ""
	}

	if reservation.GetCheckOut() != nil {
		res.CheckOut = reservation.GetCheckOut().AsTime().Format(time.RFC822)
	} else {
		res.CheckOut = ""
	}

	return res
}

func buildInquiryReservationTablesProto(tables []*cssProto.CTable) []openapi.CTable {
	response := []openapi.CTable{}
	for _, table := range tables {
		response = append(response, openapi.CTable{
			Id:          int32(table.GetId()),
			TableNumber: table.GetTableNumber(),
		})
	}

	return response

}

// Inquiry
func buildInquiryResponse(inquiry *cssProto.Inquiry) openapi.Inquiry {
	res := openapi.Inquiry{
		Id:          inquiry.GetId(),
		Guest:       buildGuestResponse(inquiry.GetGuest()),
		Branch:      buildBranchResponse(inquiry.GetBranch()),
		Status:      inquiry.GetStatus(),
		Type:        buildInquiryTypeResponse(inquiry.GetType()),
		Description: inquiry.GetDescription(),
		Date:        inquiry.GetDate().AsTime().Format(time.DateOnly),
		Time:        inquiry.GetTime().AsTime().Format("15:04"),
		Creator:     buildCreatorResponse(inquiry.GetCreator()),
		CreatedAt:   inquiry.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt:   inquiry.GetUpdatedAt().AsTime().Format(time.RFC822),
	}

	if inquiry.GetReservation() != nil {
		convReservation := buildReservationResponse(inquiry.GetReservation())
		res.Reservation = &convReservation
	} else {
		res.Reservation = nil
	}

	if inquiry.Solution != nil {
		solution := buildInquirySolutionResponse(inquiry.Solution)
		res.Solution = &solution
	} else {
		res.Solution = nil
	}

	return res
}

// All inquiries response
func buildAllInquiriesResponse(inquiries []*cssProto.Inquiry) []openapi.Inquiry {
	response := []openapi.Inquiry{}
	for _, inquiry := range inquiries {
		response = append(response, buildInquiryResponse(inquiry))
	}

	return response
}

// buildReservationInquiryStatisticsResponse
func buildReservationInquiryStatisticsResponse(reservations []*cssProto.ReservationAndInquiryStatistics) []openapi.ReservationInquiryStatistics {
	response := []openapi.ReservationInquiryStatistics{}
	for _, reservation := range reservations {
		response = append(response, openapi.ReservationInquiryStatistics{
			Hour:             reservation.GetHour().AsTime().Format("15:04"),
			ReservationCount: reservation.GetReservationCount(),
			InquiryCount:     reservation.GetInquiryCount(),
		})
	}

	return response
}

// Pagination
func buildPaginationResponse(pagination *cssProto.CPagination) openapi.Pagination {
	return openapi.Pagination{
		Total:       int32(pagination.GetTotal()),
		CurrentPage: int32(pagination.GetCurrentPage()),
		PerPage:     int32(pagination.GetPerPage()),
		From:        int32(pagination.GetFrom()),
		To:          int32(pagination.GetTo()),
		LastPage:    int32(pagination.GetLastPage()),
	}
}

func buildAllEmployeesReportsResponse(reports []*cssProto.EmployeeReport) []openapi.EmployeeReport {
	response := []openapi.EmployeeReport{}
	for _, report := range reports {
		response = append(response, buildEmployeeReportResponse(report))
	}

	return response
}

func buildEmployeeReportResponse(reports *cssProto.EmployeeReport) openapi.EmployeeReport {
	return openapi.EmployeeReport{
		Id:                reports.GetId(),
		Name:              reports.GetName(),
		Avatar:            reports.GetAvatar(),
		Role:              reports.GetRole(),
		ReservationsCount: reports.GetReservationsCount(),
		InquiriesCount:    reports.GetInquiriesCount(),
		AvgTime:           reports.GetAvgTime(),
		Rate:              reports.GetRate(),
	}
}

func buildStatisticsResponse(stats *cssProto.EmployeeStatisticsCard) openapi.EmployeeStatisticsCard {
	return openapi.EmployeeStatisticsCard{
		TotalRatings:    stats.GetTotalCount(),
		TotalAvgRatings: float32(stats.GetTotalRating()),
		Stars: openapi.StarsCard{
			FiveStar:  stats.GetStars().GetFiveStar(),
			FourStar:  stats.GetStars().GetFourStar(),
			ThreeStar: stats.GetStars().GetThreeStar(),
			TwoStar:   stats.GetStars().GetTwoStar(),
			OneStar:   stats.GetStars().GetOneStar(),
		},
	}
}

func buildAllEmployeeOperationsResponse(operations []*cssProto.EmployeeOperation) []openapi.EmployeeOperation {
	response := []openapi.EmployeeOperation{}
	for _, operation := range operations {
		response = append(response, buildEmployeeOperationResponse(operation))
	}

	return response
}

func buildEmployeeOperationResponse(operation *cssProto.EmployeeOperation) openapi.EmployeeOperation {
	return openapi.EmployeeOperation{
		Id:     operation.GetId(),
		Type:   operation.GetType(),
		Guest:  buildGuestResponse(operation.GetGuest()),
		Date:   operation.GetDate().AsTime().Format(time.DateOnly),
		Status: operation.GetStatus(),
		Timer:  operation.GetTimer(),
	}
}
