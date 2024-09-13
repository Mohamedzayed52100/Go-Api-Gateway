package guest

import (
	"time"

	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

// Guest
func buildGuestAPIResponse(proto *guestProto.Guest) openapi.Guest {
	res := openapi.Guest{
		Id:                  proto.GetId(),
		FirstName:           proto.GetFirstName(),
		LastName:            proto.GetLastName(),
		Email:               proto.GetEmail(),
		PhoneNumber:         proto.GetPhoneNumber(),
		Language:            proto.GetLanguage(),
		TotalVisits:         proto.GetTotalVisits(),
		CurrentMood:         proto.GetCurrentMood(),
		TotalSpent:          proto.GetTotalSpent(),
		TotalNoShow:         proto.GetTotalNoShow(),
		TotalCancel:         proto.GetTotalCancel(),
		UpcomingReservation: proto.GetUpcomingReservation(),
		Branches:            buildBranchVisitsResponse(proto.GetBranches()),
		Tags:                buildTagsResponse(proto.GetTags()),
		Address:             proto.GetAddress(),
		Gender:              proto.GetGender(),
		CreatedAt:           proto.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt:           proto.GetUpdatedAt().AsTime().Format(time.RFC822),
	}

	if proto.GetBirthDate() != nil {
		res.BirthDate = proto.GetBirthDate().AsTime().Format(time.DateOnly)
	} else {
		res.BirthDate = ""
	}

	if proto.GetLastVisit() != nil {
		res.LastVisit = proto.GetLastVisit().AsTime().Format(time.DateOnly)
	} else {
		res.LastVisit = ""
	}

	if res.Branches == nil {
		res.Branches = []openapi.GuestBranchVisits{}
	}

	if res.Tags == nil {
		res.Tags = []openapi.Tag{}
	}

	if proto.GetNotes() != nil {
		res.Notes = buildGuestNotesResponse(proto.GetNotes())
	} else {
		res.Notes = []openapi.GuestNote{}
	}

	return res
}

func buildGuestShortResponse(proto *guestProto.GuestShort) openapi.GuestShort {
	return openapi.GuestShort{
		Id:          proto.GetId(),
		FirstName:   proto.GetFirstName(),
		LastName:    proto.GetLastName(),
		PhoneNumber: proto.GetPhoneNumber(),
	}
}

func buildBranchVisitsResponse(proto []*guestProto.GuestBranchVisits) []openapi.GuestBranchVisits {
	var branches []openapi.GuestBranchVisits
	for _, branch := range proto {
		branches = append(branches, openapi.GuestBranchVisits{
			BranchName: branch.BranchName,
			Visits:     branch.Visits,
		})
	}

	return branches
}

func buildTagsResponse(proto []*guestProto.Tag) []openapi.Tag {
	var tags []openapi.Tag
	for _, tag := range proto {
		tg := buildTagCategoryAPIResponse(tag.Category)
		tags = append(tags, openapi.Tag{
			Id:       tag.Id,
			Category: &tg,
			Name:     tag.Name,
		})
	}

	return tags
}

func buildTagCategoryAPIResponse(proto *guestProto.TagCategory) openapi.TagCategory {
	return openapi.TagCategory{
		Id:             proto.Id,
		Name:           proto.Name,
		Color:          proto.Color,
		Classification: proto.Classification,
		OrderIndex:     proto.OrderIndex,
	}
}

func buildGuestNotesResponse(proto []*guestProto.GuestNote) []openapi.GuestNote {
	var notes []openapi.GuestNote
	for _, note := range proto {
		notes = append(notes, buildGuestNoteAPIResponse(note))
	}

	return notes
}

func buildGuestNoteAPIResponse(proto *guestProto.GuestNote) openapi.GuestNote {
	return openapi.GuestNote{
		Id:          proto.Id,
		GuestId:     proto.GuestId,
		Description: proto.Description,
		Creator:     buildCreatorAPIResponse(proto.Creator),
		CreatedAt:   proto.CreatedAt.AsTime().Format(time.RFC822),
		UpdatedAt:   proto.UpdatedAt.AsTime().Format(time.RFC822),
	}
}

func buildCreatorAPIResponse(proto *guestProto.CreatorProfile) openapi.User {
	return openapi.User{
		Id:          proto.GetId(),
		FirstName:   proto.GetFirstName(),
		LastName:    proto.GetLastName(),
		PhoneNumber: proto.GetPhoneNumber(),
		Avatar:      proto.GetAvatar(),
		Email:       proto.GetEmail(),
		Role:        proto.GetRole(),
	}
}

func buildGuestsAPIResponse(proto []*guestProto.Guest) []openapi.Guest {
	result := []openapi.Guest{}
	for _, guest := range proto {
		result = append(result, buildGuestAPIResponse(guest))
	}
	return result
}

// Seating Area
func buildSeatingAreaResponse(area *guestProto.SeatingArea) openapi.SeatingArea {
	return openapi.SeatingArea{
		Id:   area.GetId(),
		Name: area.GetName(),
	}
}

// Table
func buildTablesResponse(table []*guestProto.Table) []openapi.Table {
	var result []openapi.Table
	for _, t := range table {
		result = append(result, openapi.Table{
			Id:           t.GetId(),
			SeatingArea:  buildSeatingAreaResponse(t.GetSeatingArea()),
			TableNumber:  t.GetTableNumber(),
			PosNumber:    t.GetPosNumber(),
			MinPartySize: t.GetMinPartySize(),
			MaxPartySize: t.GetMaxPartySize(),
			CreatedAt:    t.GetCreatedAt().AsTime().Format(time.RFC822),
			UpdatedAt:    t.GetUpdatedAt().AsTime().Format(time.RFC822),
		})
	}
	return result
}

// Reservation Status
func buildReservationStatusResponse(status *guestProto.ReservationStatus) openapi.ReservationStatus {
	return openapi.ReservationStatus{
		Id:        status.Id,
		Name:      status.Name,
		Category:  status.Category,
		Type:      status.Type,
		Color:     status.Color,
		Icon:      status.Icon,
		CreatedAt: status.CreatedAt.AsTime().Format(time.RFC822),
		UpdatedAt: status.UpdatedAt.AsTime().Format(time.RFC822),
	}
}

// Reservation Feedback
func buildAllReservationFeedbacksResponse(feedback []*guestProto.ReservationFeedback) []openapi.ReservationFeedback {
	result := []openapi.ReservationFeedback{}

	for _, f := range feedback {
		result = append(result, buildReservationFeedbackAPIResponse(f))
	}

	return result
}

func buildAllFeedbackSectionsResponse(sections []*guestProto.FeedbackSection) []openapi.FeedbackSection {
	result := []openapi.FeedbackSection{}
	for _, s := range sections {
		result = append(result, buildFeedbackSectionResponse(s))
	}
	return result
}

func buildFeedbackSectionResponse(section *guestProto.FeedbackSection) openapi.FeedbackSection {
	return openapi.FeedbackSection{
		Id:        section.Id,
		Name:      section.Name,
		CreatedAt: section.CreatedAt.AsTime().Format(time.RFC822),
		UpdatedAt: section.UpdatedAt.AsTime().Format(time.RFC822),
	}
}

func buildReservationFeedbackAPIResponse(feedback *guestProto.ReservationFeedback) openapi.ReservationFeedback {
	rate := feedback.GetRate()
	description := feedback.GetDescription()
	res := openapi.ReservationFeedback{
		Id:          feedback.GetId(),
		Reservation: buildReservationShortResponse(feedback.GetReservation()),
		Guest:       buildGuestShortResponse(feedback.GetGuest()),
		Sections:    buildAllFeedbackSectionsResponse(feedback.GetSections()),
		Status:      feedback.GetStatus(),
		Rate:        &rate,
		Description: &description,
		CreatedAt:   feedback.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt:   feedback.GetUpdatedAt().AsTime().Format(time.RFC822),
	}

	if feedback.GetSolution() != nil {
		solution := buildReservationFeedbackSolutionResponse(feedback.GetSolution())
		res.Solution = &solution
	} else {
		res.Solution = nil
	}

	if res.Sections == nil {
		res.Sections = []openapi.FeedbackSection{}
	} else {
		res.Sections = buildAllFeedbackSectionsResponse(feedback.GetSections())
	}

	if feedback.Reservation.Tables == nil {
		res.Reservation.Tables = []openapi.Table{}
	}

	return res
}

// Short Reservation Feedback with guest and reservation
func buildShortReservationFeedbackAPIResponse(feedback *guestProto.ReservationFeedbackShort) openapi.ReservationFeedbackShort {
	description := feedback.GetDescription()
	if description == "" {
		description = ""
	}

	rate := feedback.GetRate()
	if rate == 0 {
		rate = 0
	}

	return openapi.ReservationFeedbackShort{
		Id:          feedback.GetId(),
		Rate:        feedback.GetRate(),
		Description: &description,
		CreatedAt:   feedback.GetCreatedAt().AsTime().Format(time.RFC822),
	}
}

// Reservation Special Occasion
func buildReservationSpecialOccasionResponse(occasion *guestProto.ReservationSpecialOccasion) openapi.SpecialOccasion {
	return openapi.SpecialOccasion{
		Id:        occasion.Id,
		Name:      occasion.Name,
		Color:     occasion.Color,
		Icon:      occasion.Icon,
		CreatedAt: occasion.CreatedAt.AsTime().Format(time.RFC822),
		UpdatedAt: occasion.UpdatedAt.AsTime().Format(time.RFC822),
	}
}

// Order
func buildReservationOrderResponse(order *guestProto.ReservationOrder) *openapi.ReservationOrder {
	res := &openapi.ReservationOrder{
		Id:             order.GetId(),
		DiscountAmount: order.GetDiscountAmount(),
		DiscountReason: order.GetDiscountReason(),
		PrevailingTax:  order.GetPrevailingTax(),
		Tax:            order.GetTax(),
		SubTotal:       order.GetSubTotal(),
		FinalTotal:     order.GetFinalTotal(),
		CreatedAt:      order.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt:      order.GetUpdatedAt().AsTime().Format(time.RFC822),
	}

	if order.GetItems() != nil {
		res.Items = buildReservationOrderItemsResponse(order.GetItems())
	} else {
		res.Items = []openapi.ReservationOrderItem{}
	}

	return res
}

func buildReservationOrderItemsResponse(items []*guestProto.ReservationOrderItem) []openapi.ReservationOrderItem {
	var result []openapi.ReservationOrderItem
	for _, item := range items {
		result = append(result, openapi.ReservationOrderItem{
			Id:        item.GetId(),
			ItemName:  item.GetItemName(),
			Quantity:  item.GetQuantity(),
			Cost:      item.GetCost(),
			CreatedAt: item.GetCreatedAt().AsTime().Format(time.RFC822),
			UpdatedAt: item.GetUpdatedAt().AsTime().Format(time.RFC822),
		})
	}
	return result

}

func buildReservationGuestResponse(proto *guestProto.ReservationGuest) openapi.ReservationGuest {
	res := openapi.ReservationGuest{
		Id:          proto.GetId(),
		FirstName:   proto.GetFirstName(),
		LastName:    proto.GetLastName(),
		PhoneNumber: proto.GetPhoneNumber(),
		TotalVisits: proto.GetTotalVisits(),
		TotalSpent:  proto.GetTotalSpent(),
		TotalNoShow: proto.GetTotalNoShow(),
		TotalCancel: proto.GetTotalCancel(),
		IsPrimary:   proto.GetIsPrimary(),
		Gender:      proto.GetGender(),
		Tags:        buildTagsResponse(proto.GetTags()),
	}

	if proto.GetNote() != nil {
		note := buildGuestNoteAPIResponse(proto.GetNote())
		res.Note = &note
	} else if proto.GetIsPrimary() {
		res.Note = &openapi.GuestNote{}
	}

	if len(proto.GetTags()) == 0 {
		res.Tags = []openapi.Tag{}
	}

	return res
}

func buildAllReservationGuestsResponse(proto []*guestProto.ReservationGuest) []openapi.ReservationGuest {
	var result []openapi.ReservationGuest
	for _, guest := range proto {
		result = append(result, buildReservationGuestResponse(guest))
	}
	return result
}

// Reservation
func buildReservationAPIResponse(proto *guestProto.Reservation) openapi.Reservation {
	res := openapi.Reservation{
		Id:             proto.GetId(),
		ReservationRef: proto.GetReservationRef(),
		Guests:         buildAllReservationGuestsResponse(proto.GetGuests()),
		Branch: openapi.ReservationBranch{
			Id:   proto.GetBranch().GetId(),
			Name: proto.GetBranch().GetName(),
		},
		Shift: openapi.ReservationShift{
			Id:   proto.GetShift().GetId(),
			Name: proto.GetShift().GetName(),
			Cast: openapi.Cast{},
		},
		Tables:       buildTablesResponse(proto.GetTables()),
		GuestsNumber: proto.GetGuestsNumber(),
		SeatedGuests: proto.GetSeatedGuests(),
		Date:         proto.GetDate().AsTime().Format(time.DateOnly),
		Time:         proto.GetTime().AsTime().Format("15:04"),
		ReservedVia:  proto.GetReservedVia(),
		Status:       buildReservationStatusResponse(proto.GetStatus()),
		TotalSpent:   proto.GetTotalSpent(),
		Payment: openapi.ReservationPayment{
			Status:      proto.GetPayment().GetStatus(),
			TotalPaid:   proto.GetPayment().GetTotalPaid(),
			TotalUnPaid: proto.GetPayment().GetTotalUnPaid(),
		},
		CreatedAt: proto.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt: proto.GetUpdatedAt().AsTime().Format(time.RFC822),
	}

	if proto.GetCheckIn() != nil {
		res.CheckIn = proto.GetCheckIn().AsTime().Format(time.RFC822)
	}

	if proto.GetCheckOut() != nil {
		res.CheckOut = proto.GetCheckOut().AsTime().Format(time.RFC822)
	}

	if proto.GetFeedback() != nil {
		feedback := buildShortReservationFeedbackAPIResponse(proto.GetFeedback())
		res.Feedback = &feedback
	} else {
		res.Feedback = &openapi.ReservationFeedbackShort{}
	}

	if proto.GetSpecialOccasion() != nil {
		sp := buildReservationSpecialOccasionResponse(proto.GetSpecialOccasion())
		res.SpecialOccasion = &sp
	} else {
		res.SpecialOccasion = nil
	}

	if proto.GetTags() != nil {
		res.Tags = buildTagsResponse(proto.GetTags())
	} else {
		res.Tags = []openapi.Tag{}
	}

	if proto.GetNote() != nil {
		res.Note = buildReservationNoteAPIResponse(proto.GetNote())
	} else {
		res.Note = openapi.ReservationNote{}
	}

	if proto.GetCreator() != nil {
		creator := buildCreatorAPIResponse(proto.GetCreator())
		res.Creator = &creator
	} else {
		res.Creator = nil
	}

	if proto.GetTables() != nil {
		res.Tables = buildTablesResponse(proto.GetTables())
	} else {
		res.Tables = []openapi.Table{}
	}

	return res
}

func buildReservationShortResponse(proto *guestProto.ReservationShort) openapi.ReservationShort {
	res := openapi.ReservationShort{
		Id:           proto.GetId(),
		GuestsNumber: proto.GetGuestsNumber(),
		SeatedGuests: proto.GetSeatedGuests(),
		Date:         proto.GetDate().AsTime().Format(time.DateOnly),
		Time:         proto.GetTime().AsTime().Format("15:04"),
		ReservedVia:  proto.GetReservedVia(),
		Status:       buildReservationStatusResponse(proto.GetStatus()),
		Tables:       buildTablesResponse(proto.GetTables()),
		CreatedAt:    proto.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt:    proto.GetUpdatedAt().AsTime().Format(time.RFC822),
	}

	if proto.GetCheckIn() != nil {
		res.CheckIn = proto.GetCheckIn().AsTime().Format(time.RFC822)
	}

	if proto.GetCheckOut() != nil {
		res.CheckOut = proto.GetCheckOut().AsTime().Format(time.RFC822)
	}

	if proto.GetSpecialOccasion() != nil {
		sp := buildReservationSpecialOccasionResponse(proto.GetSpecialOccasion())
		res.SpecialOccasion = &sp
	} else {
		res.SpecialOccasion = nil
	}

	if proto.GetBranch() != nil {
		res.Branch = openapi.CBranch{
			Id:   proto.GetBranch().GetId(),
			Name: proto.GetBranch().GetName(),
		}
	}

	return res
}

func buildReservationsAPIResponse(proto []*guestProto.Reservation) []openapi.Reservation {
	result := []openapi.Reservation{}
	for _, reservation := range proto {
		result = append(result, buildReservationAPIResponse(reservation))
	}
	return result
}

// Available Times
func buildAvailableTimesResponse(times []*guestProto.AvailableTime) []openapi.AvailableTime {
	result := []openapi.AvailableTime{}
	for _, t := range times {
		result = append(result, openapi.AvailableTime{
			Time:      t.GetTime().AsTime().Format("15:04"),
			Pacing:    t.GetPacing(),
			Capacity:  t.GetCapacity(),
			Available: t.GetAvailable(),
		})
	}

	return result
}

// Logs
func buildReservationLogResponse(log *guestProto.ReservationLog) openapi.ReservationLog {
	res := openapi.ReservationLog{
		Id:            log.GetId(),
		ReservationId: log.GetReservationId(),
		MadeBy:        log.GetMadeBy(),
		FieldName:     log.GetFieldName(),
		OldValue:      log.GetOldValue(),
		NewValue:      log.GetNewValue(),
		Action:        log.GetAction(),
		CreatedAt:     log.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt:     log.GetUpdatedAt().AsTime().Format(time.RFC822),
	}

	if log.GetCreator() != nil {
		creator := buildCreatorAPIResponse(log.GetCreator())
		res.Creator = &creator
	} else {
		res.Creator = nil
	}

	return res
}

func buildAllReservationLogsResponse(logs []*guestProto.ReservationLog) []openapi.ReservationLog {
	result := make([]openapi.ReservationLog, 0)
	for _, log := range logs {
		result = append(result, buildReservationLogResponse(log))
	}

	return result
}

func buildReservationWaitlistLogResponse(log *guestProto.ReservationWaitlistLog) openapi.ReservationWaitlistLog {
	res := openapi.ReservationWaitlistLog{
		Id:                    log.GetId(),
		ReservationWaitlistId: log.GetReservationWaitlistId(),
		MadeBy:                log.GetMadeBy(),
		FieldName:             log.GetFieldName(),
		OldValue:              log.GetOldValue(),
		NewValue:              log.GetNewValue(),
		Action:                log.GetAction(),
		CreatedAt:             log.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt:             log.GetUpdatedAt().AsTime().Format(time.RFC822),
	}

	if log.GetCreator() != nil {
		creator := buildCreatorAPIResponse(log.GetCreator())
		res.Creator = &creator
	} else {
		res.Creator = nil
	}

	return res
}

func buildAllReservationWaitlistLogsResponse(logs []*guestProto.ReservationWaitlistLog) []openapi.ReservationWaitlistLog {
	result := make([]openapi.ReservationWaitlistLog, 0)
	for _, log := range logs {
		result = append(result, buildReservationWaitlistLogResponse(log))
	}

	return result
}

func buildGuestLogResponse(log *guestProto.GuestLog) openapi.GuestLog {
	res := openapi.GuestLog{
		Id:        log.GetId(),
		GuestId:   log.GetGuestId(),
		MadeBy:    log.GetMadeBy(),
		FieldName: log.GetFieldName(),
		OldValue:  log.GetOldValue(),
		NewValue:  log.GetNewValue(),
		Action:    log.GetAction(),
		CreatedAt: log.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt: log.GetUpdatedAt().AsTime().Format(time.RFC822),
	}

	if log.GetCreator() != nil {
		creator := buildCreatorAPIResponse(log.GetCreator())
		res.Creator = &creator
	} else {
		res.Creator = nil
	}

	return res
}

func buildAllGuestLogsResponse(logs []*guestProto.GuestLog) []openapi.GuestLog {
	result := make([]openapi.GuestLog, 0)
	for _, log := range logs {
		result = append(result, buildGuestLogResponse(log))
	}

	return result
}

func buildReservationNoteAPIResponse(note *guestProto.ReservationNote) openapi.ReservationNote {
	res := openapi.ReservationNote{
		Id:          note.Id,
		Description: note.Description,
		CreatedAt:   note.CreatedAt.AsTime().Format(time.RFC822),
		UpdatedAt:   note.UpdatedAt.AsTime().Format(time.RFC822),
	}

	if note.Creator == nil {
		res.Creator = &openapi.User{}
	} else {
		creator := buildCreatorAPIResponse(note.GetCreator())
		res.Creator = &creator
	}

	return res
}

func buildWaitingReservationNoteResponse(note *guestProto.ReservationWaitlistNote) openapi.ReservationWaitlistNote {
	res := openapi.ReservationWaitlistNote{
		Id:                    note.GetId(),
		ReservationWaitlistId: note.GetReservationWaitlistId(),
		Description:           note.GetDescription(),
		CreatedAt:             note.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt:             note.GetUpdatedAt().AsTime().Format(time.RFC822),
	}

	if note.GetCreator() != nil {
		creator := buildCreatorAPIResponse(note.GetCreator())
		res.Creator = &creator
	}

	return res
}

func buildReservationWaitlistResponse(res *guestProto.ReservationWaitlist) openapi.ReservationWaitlist {
	result := openapi.ReservationWaitlist{
		Id:           res.GetId(),
		Guest:        buildGuestAPIResponse(res.GetGuest()),
		SeatingArea:  buildSeatingAreaResponse(res.GetSeatingArea()),
		GuestsNumber: res.GetGuestsNumber(),
		WaitingTime:  res.GetWaitingTime(),
		CreatedAt:    res.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt:    res.GetUpdatedAt().AsTime().Format(time.RFC822),
	}

	if res.GetNote() != nil {
		result.Note = buildWaitingReservationNoteResponse(res.GetNote())
	} else {
		result.Note = openapi.ReservationWaitlistNote{}
	}

	if res.GetTags() != nil {
		result.Tags = buildTagsResponse(res.GetTags())
	} else {
		result.Tags = []openapi.Tag{}
	}

	return result
}

// Build wait list response
func buildAllReservationWaitlistsResponse(waitlist []*guestProto.ReservationWaitlist) []openapi.ReservationWaitlist {
	result := make([]openapi.ReservationWaitlist, 0)
	for _, w := range waitlist {
		result = append(result, buildReservationWaitlistResponse(w))
	}
	return result
}

// Guest statistics response
func buildGuestStatisticsResponse(proto *guestProto.GuestStatistics) openapi.GuestStatistics {
	return openapi.GuestStatistics{
		TotalReservations:  proto.GetTotalReservations(),
		TotalSpent:         proto.GetTotalSpent(),
		PublicSatisfaction: proto.GetPublicSatisfaction(),
	}
}

// Months spending response
func buildGuestSpendingMonthsResponse(proto []*guestProto.MonthSpending) []openapi.MonthSpending {
	result := []openapi.MonthSpending{}
	for _, month := range proto {
		result = append(result, openapi.MonthSpending{
			Month:      month.GetMonth(),
			TotalSpent: month.GetTotalSpent(),
		})
	}

	return result
}

// Guest spending response
func buildGuestSpendingResponse(proto []*guestProto.YearSpending) []openapi.YearSpending {
	result := []openapi.YearSpending{}
	for _, spending := range proto {
		result = append(result, openapi.YearSpending{
			Year:   spending.GetYear(),
			Months: buildGuestSpendingMonthsResponse(spending.GetMonths()),
		})
	}

	return result
}

// Guest reservation statistics response
func buildGuestReservationStatisticsResponse(proto []*guestProto.GuestReservationStatistics) []openapi.GuestReservationStatistics {
	result := []openapi.GuestReservationStatistics{}
	for _, statistics := range proto {
		result = append(result, openapi.GuestReservationStatistics{
			Name:  statistics.GetName(),
			Value: statistics.GetValue(),
		})
	}

	return result
}

// All special occasions response
func buildAllSpecialOccasionsResponse(proto []*guestProto.SpecialOccasion) []openapi.SpecialOccasion {
	result := []openapi.SpecialOccasion{}
	for _, special := range proto {
		result = append(result, openapi.SpecialOccasion{
			Id:        special.GetId(),
			Name:      special.GetName(),
			Color:     special.GetColor(),
			Icon:      special.GetIcon(),
			CreatedAt: special.GetCreatedAt().AsTime().Format(time.RFC822),
			UpdatedAt: special.GetUpdatedAt().AsTime().Format(time.RFC822),
		})
	}

	return result
}

func buildCreatorResponse(creator *guestProto.CreatorProfile) openapi.Creator {
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

func buildReservationFeedbackCommentResponse(comment *guestProto.ReservationFeedbackComment) openapi.ReservationFeedbackComment {
	res := openapi.ReservationFeedbackComment{
		Id:        int32(comment.GetId()),
		Creator:   buildCreatorResponse(comment.GetCreator()),
		Comment:   comment.GetComment(),
		CreatedAt: comment.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt: comment.GetUpdatedAt().AsTime().Format(time.RFC822),
	}

	return res
}

func buildAllReservationFeedbackCommentsResponse(comments []*guestProto.ReservationFeedbackComment) []openapi.ReservationFeedbackComment {
	response := []openapi.ReservationFeedbackComment{}
	for _, comment := range comments {
		response = append(response, buildReservationFeedbackCommentResponse(comment))
	}

	return response
}

func buildReservationFeedbackSolutionResponse(solution *guestProto.ReservationFeedbackSolution) openapi.ReservationFeedbackSolution {
	return openapi.ReservationFeedbackSolution{
		Id:        int32(solution.GetId()),
		Creator:   buildCreatorResponse(solution.Creator),
		Solution:  solution.GetSolution(),
		CreatedAt: solution.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt: solution.GetUpdatedAt().AsTime().Format(time.RFC822),
	}
}

func buildCoverFlowResponse(coverFlow []*guestProto.CoverFlow) []openapi.CoverFlow {
	var response []openapi.CoverFlow
	for _, c := range coverFlow {
		response = append(response, openapi.CoverFlow{
			Time:         c.GetTime(),
			Reservations: buildCoverFlowReservationsResponse(c.GetReservations()),
		})
	}

	return response
}

func buildCoverFlowReservationsResponse(reservations []*guestProto.CoverFlowReservation) []openapi.CoverFlowReservation {
	response := []openapi.CoverFlowReservation{}
	for _, r := range reservations {
		response = append(response, openapi.CoverFlowReservation{
			Id:           r.GetId(),
			GuestsNumber: r.GetGuestsNumber(),
			Status: openapi.CoverFlowReservationStatus{
				Id:    r.GetStatus().GetId(),
				Name:  r.GetStatus().GetName(),
				Color: r.GetStatus().GetColor(),
				Icon:  r.GetStatus().GetIcon(),
			},
		})
	}
	return response
}

func buildWidgetAvailableTimesResponse(times []*guestProto.AvailableTime) []openapi.WidgetAvailableTime {
	result := []openapi.WidgetAvailableTime{}
	for _, t := range times {
		result = append(result, openapi.WidgetAvailableTime{
			Time:      t.GetTime().AsTime().Format("15:04"),
			Available: t.GetAvailable(),
		})
	}

	return result
}

func buildAllPaymentsResponse(payments []*guestProto.PaymentResponse) []openapi.PaymentResponse {
	result := []openapi.PaymentResponse{}
	for _, payment := range payments {
		result = append(result, buildPaymentResponse(payment))
	}

	return result
}

func buildPaymentResponse(payment *guestProto.PaymentResponse) openapi.PaymentResponse {
	res := openapi.PaymentResponse{
		Id:       payment.Id,
		Contacts: payment.Contacts,
		Status:   payment.Status,
		Guest: openapi.CGuest{
			FirstName:   payment.Guest.FirstName,
			LastName:    payment.Guest.LastName,
			PhoneNumber: payment.Guest.PhoneNumber,
			Address:     payment.Guest.Address,
		},
		Delivery: payment.Delivery,
		Invoice: openapi.Invoice{
			InvoiceId:  payment.Invoice.InvoiceId,
			InvoiceRef: payment.Invoice.InvoiceRef,
			Date:       payment.Invoice.Date,
			Waiter:     "",
			Items:      buildPaymentItemsResponse(payment.Invoice.Items),
			SubTotal:   payment.Invoice.SubTotal,
		},
		Branch: openapi.PaymentBranch{
			Name:          payment.Branch.Name,
			Address:       payment.Branch.Address,
			VatPercent:    payment.Branch.VatPercent,
			ServiceCharge: payment.Branch.ServiceCharge,
			CrNumber:      payment.Branch.CrNumber,
			VatRegNumber:  payment.Branch.VatRegNumber,
		},
	}

	if payment.Status == "paid" {
		res.Card = &openapi.PaymentCard{
			LastFourDigits: payment.Card.LastFourDigits,
			CardType:       payment.Card.CardType,
			CardExpireDate: payment.Card.CardExpireDate,
		}
	} else {
		res.Card = nil
	}

	return res
}

func buildPaymentItemsResponse(items []*guestProto.PaymentItem) []openapi.PaymentItem {
	var result []openapi.PaymentItem
	for _, item := range items {
		result = append(result, openapi.PaymentItem{
			Id:       item.Id,
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
		})
	}
	return result
}
