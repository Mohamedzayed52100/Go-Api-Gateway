package settings

import (
	"log"
	"strings"
	"time"

	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

// Tag Category
func buildTagCategoryResponse(category *settingsProto.TagCategory) openapi.TagCategory {
	return openapi.TagCategory{
		Id:             category.GetId(),
		Name:           category.GetName(),
		Color:          category.GetColor(),
		Classification: category.GetClassification(),
		OrderIndex:     category.GetOrderIndex(),
		Tags:           buildCategoryTagsWithoutCategoryResponse(category.GetTags()),
		IsDisabled:     category.GetIsDisabled(),
	}
}

func buildAllTagCategoriesResponse(categories []*settingsProto.TagCategory) []openapi.TagCategory {
	var result []openapi.TagCategory
	for _, category := range categories {
		result = append(result, buildTagCategoryResponse(category))
	}

	return result
}

func buildCategoryTagsWithoutCategoryResponse(tags []*settingsProto.Tag) []openapi.Tag {
	var result []openapi.Tag
	for _, tag := range tags {
		result = append(result, openapi.Tag{
			Id:        tag.GetId(),
			Name:      tag.GetName(),
			CreatedAt: tag.GetCreatedAt().AsTime().Format(time.RFC822),
			UpdatedAt: tag.GetUpdatedAt().AsTime().Format(time.RFC822),
		})
	}
	return result
}

func buildTagResponse(tag *settingsProto.Tag) openapi.Tag {
	res := openapi.Tag{
		Id:        tag.GetId(),
		Name:      tag.GetName(),
		CreatedAt: tag.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt: tag.GetUpdatedAt().AsTime().Format(time.RFC822),
	}

	if tag.GetCategory() != nil {
		tg := buildTagCategoryResponse(tag.GetCategory())
		res.Category = &tg
	}

	return res
}

func buildAllTagsResponse(tags []*settingsProto.Tag) []openapi.Tag {
	var result []openapi.Tag

	for _, tag := range tags {
		result = append(result, buildTagResponse(tag))
	}

	return result
}

// Reservation Status
func buildReservationStatusResponse(status *settingsProto.ReservationStatus) openapi.ReservationStatus {
	res := openapi.ReservationStatus{
		Id:       status.GetId(),
		Name:     status.GetName(),
		Category: status.GetCategory(),
		Type:     status.GetType(),
		Color:    status.GetColor(),
		Icon:     status.GetIcon(),
	}

	if status.GetCreatedAt() != nil {
		res.CreatedAt = status.CreatedAt.AsTime().Format(time.RFC822)
	}
	if status.GetUpdatedAt() != nil {
		res.UpdatedAt = status.UpdatedAt.AsTime().Format(time.RFC822)
	}

	return res
}
func buildAllReservationStatusesResponse(statuses []*settingsProto.ReservationStatus) []openapi.ReservationStatus {
	var result []openapi.ReservationStatus
	for _, status := range statuses {
		result = append(result, buildReservationStatusResponse(status))
	}
	return result
}

func buildAllBranchReservationStatusesResponse(proto []*settingsProto.BranchReservationStatus) []openapi.BranchReservationStatus {
	result := []openapi.BranchReservationStatus{}
	for _, item := range proto {
		result = append(result, openapi.BranchReservationStatus{
			Branch:   item.Branch,
			Statuses: buildAllReservationStatusesResponse(item.Statuses),
		})
	}

	return result
}

// Pacing response
func buildPacingResponse(pacings []*settingsProto.Pacing) []openapi.Pacing {
	result := make([]openapi.Pacing, 0)
	for _, pacing := range pacings {
		result = append(result, openapi.Pacing{
			Id:            pacing.GetId(),
			ShiftId:       pacing.GetShiftId(),
			SeatingAreaId: pacing.GetSeatingAreaId(),
			Hour:          pacing.GetHour(),
			Capacity:      pacing.GetCapacity(),
			CreatedAt:     pacing.GetCreatedAt().AsTime().Format(time.RFC822),
			UpdatedAt:     pacing.GetUpdatedAt().AsTime().Format(time.RFC822),
		})
	}

	return result
}

// Shift
func buildShiftResponse(shift *settingsProto.Shift) openapi.Shift {
	res := openapi.Shift{
		Id:           shift.GetId(),
		Name:         shift.GetName(),
		From:         shift.GetFrom().AsTime().Format("15:04"),
		To:           shift.GetTo().AsTime().Format("15:04"),
		StartDate:    shift.GetStartDate().AsTime().Format(time.DateOnly),
		EndDate:      shift.GetEndDate().AsTime().Format(time.DateOnly),
		TimeInterval: shift.GetTimeInterval(),
		FloorPlan:    buildShortFloorPlanResponse(shift.GetFloorPlan()),
		SeatingAreas: buildAllSeatingAreasResponse(shift.GetSeatingAreas()),
		Category:     buildShiftCategoryResponse(shift.GetCategory()),
		MinGuests:    shift.GetMinGuests(),
		MaxGuests:    shift.GetMaxGuests(),
		Turnover:     buildTurnoverResponse(shift.GetTurnover()),
		Pacing:       buildPacingResponse(shift.GetPacing()),
		CreatedAt:    shift.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt:    shift.GetUpdatedAt().AsTime().Format(time.RFC822),
	}

	if shift.DaysToRepeat != nil {
		res.DaysToRepeat = shift.DaysToRepeat
	} else {
		res.DaysToRepeat = make([]string, 0)
	}

	if shift.Exceptions != nil {
		res.Exceptions = shift.Exceptions
	} else {
		res.Exceptions = make([]string, 0)
	}

	return res
}

func buildShortShiftResponse(shift *settingsProto.Shift) openapi.Shift {
	res := openapi.Shift{
		Id:           shift.GetId(),
		Name:         shift.GetName(),
		From:         shift.GetFrom().AsTime().Format("15:04"),
		To:           shift.GetTo().AsTime().Format("15:04"),
		StartDate:    shift.GetStartDate().AsTime().Format(time.DateOnly),
		EndDate:      shift.GetEndDate().AsTime().Format(time.DateOnly),
		TimeInterval: shift.GetTimeInterval(),
		FloorPlan:    buildShortFloorPlanResponse(shift.GetFloorPlan()),
		SeatingAreas: buildAllSeatingAreasResponse(shift.GetSeatingAreas()),
		Category:     buildShiftCategoryResponse(shift.GetCategory()),
		MinGuests:    shift.GetMinGuests(),
		MaxGuests:    shift.GetMaxGuests(),
		Turnover:     buildTurnoverResponse(shift.GetTurnover()),
		Pacing:       buildPacingResponse(shift.GetPacing()),
		CreatedAt:    shift.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt:    shift.GetUpdatedAt().AsTime().Format(time.RFC822),
	}

	if shift.DaysToRepeat != nil {
		res.DaysToRepeat = shift.DaysToRepeat
	} else {
		res.DaysToRepeat = make([]string, 0)
	}

	if shift.Exceptions != nil {
		res.Exceptions = shift.Exceptions
	} else {
		res.Exceptions = make([]string, 0)
	}

	return res
}

func buildAllShiftsResponse(shifts []*settingsProto.Shift) []openapi.Shift {
	result := make([]openapi.Shift, 0)
	for _, shift := range shifts {
		result = append(result, buildShortShiftResponse(shift))
	}

	return result
}

func buildAllShiftCategoriesResponse(categories []*settingsProto.ShiftCategory) []openapi.ShiftCategory {
	result := make([]openapi.ShiftCategory, 0)
	for _, category := range categories {
		result = append(result, buildShiftCategoryResponse(category))
	}

	return result
}

func buildShiftCategoryResponse(category *settingsProto.ShiftCategory) openapi.ShiftCategory {
	return openapi.ShiftCategory{
		Id:        category.GetId(),
		Name:      category.GetName(),
		Color:     category.GetColor(),
		CreatedAt: category.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt: category.GetUpdatedAt().AsTime().Format(time.RFC822),
	}
}

// Element properties response
func buildElementProperty(properties *settingsProto.FloorPlanElementProperty) openapi.ElementProperty {
	return openapi.ElementProperty{
		Shape:  properties.GetShape(),
		Width:  properties.GetWidth(),
		Height: properties.GetHeight(),
		Angle:  properties.GetAngle(),
		X:      properties.GetX(),
		Y:      properties.GetY(),
	}
}

// Floor elements response
func buildFloorPlanElementsResponse(elements []*settingsProto.FloorPlanElement) []openapi.FloorPlanElement {
	result := make([]openapi.FloorPlanElement, 0)
	for _, element := range elements {
		ele := openapi.FloorPlanElement{
			Id:          element.GetId(),
			ElementType: element.GetElementType(),
			TableId:     element.GetTableId(),
			Property:    buildElementProperty(element.GetProperty()),
			CreatedAt:   element.CreatedAt.AsTime().Format(time.RFC822),
			UpdatedAt:   element.UpdatedAt.AsTime().Format(time.RFC822),
		}

		if element.GetTable() != nil {
			table := buildTableResponse(element.GetTable())
			ele.Table = &table
		} else {
			ele.Table = nil
		}

		result = append(result, ele)
	}

	return result
}

func buildFloorPlanLayoutsResponse(layouts []*settingsProto.FloorPlanLayout) []openapi.FloorPlanLayout {
	result := make([]openapi.FloorPlanLayout, 0)
	for _, layout := range layouts {
		result = append(result, buildFloorPlanLayoutResponse(layout))
	}

	return result
}

func buildFloorPlanLayoutResponse(layout *settingsProto.FloorPlanLayout) openapi.FloorPlanLayout {
	return openapi.FloorPlanLayout{
		SeatingArea: buildSeatingAreaResponse(layout.SeatingArea),
		Scale:       layout.Scale,
		Elements:    buildFloorPlanElementsResponse(layout.Elements),
	}
}

func buildFloorPlanResponse(floorPlan *settingsProto.FloorPlan) openapi.FloorPlan {
	return openapi.FloorPlan{
		Id:        floorPlan.GetId(),
		Name:      floorPlan.GetName(),
		Layouts:   buildFloorPlanLayoutsResponse(floorPlan.GetLayouts()),
		CreatedAt: floorPlan.CreatedAt.AsTime().Format(time.RFC822),
		UpdatedAt: floorPlan.UpdatedAt.AsTime().Format(time.RFC822),
	}
}

func buildShortFloorPlanResponse(floorPlan *settingsProto.ShortFloorPlan) openapi.ShortFloorPlan {
	return openapi.ShortFloorPlan{
		Id:   floorPlan.GetId(),
		Name: floorPlan.GetName(),
	}
}

func buildAllFloorPlansResponse(floorPlans []*settingsProto.FloorPlan) []openapi.FloorPlan {
	result := make([]openapi.FloorPlan, 0)
	for _, floorPlan := range floorPlans {
		result = append(result, buildFloorPlanResponse(floorPlan))
	}

	return result
}

func buildCastResponse(cast *settingsProto.Cast) openapi.Cast {
	return openapi.Cast{
		Id:        int32(cast.GetId()),
		Staff:     buildStaffResponse(cast.GetStaff()),
		CreatedAt: cast.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt: cast.GetUpdatedAt().AsTime().Format(time.RFC822),
	}
}

func buildStaffResponse(staff []*settingsProto.Staff) []openapi.Staff {
	result := make([]openapi.Staff, 0)
	for _, s := range staff {
		result = append(result, openapi.Staff{
			Id:          s.GetId(),
			CastId:      s.GetCastId(),
			Name:        s.Name,
			Role:        s.Role,
			PhoneNumber: s.PhoneNumber,
			CreatedAt:   s.GetCreatedAt().AsTime().Format(time.RFC822),
			UpdatedAt:   s.GetUpdatedAt().AsTime().Format(time.RFC822),
		})
	}
	return result
}

func buildSeatingAreaResponse(area *settingsProto.SeatingArea) openapi.SeatingArea {
	return openapi.SeatingArea{
		Id:   area.GetId(),
		Name: area.GetName(),
	}
}

func buildTurnoverResponse(turnover []*settingsProto.Turnover) []openapi.Turnover {
	result := make([]openapi.Turnover, 0)
	for _, t := range turnover {
		result = append(result, openapi.Turnover{
			Id:           t.GetId(),
			ShiftId:      t.GetShiftId(),
			GuestsNumber: t.GetGuestsNumber(),
			TurnoverTime: t.GetTurnoverTime().AsTime().Format("15:04"),
			CreatedAt:    t.GetCreatedAt().AsTime().Format(time.RFC822),
			UpdatedAt:    t.GetUpdatedAt().AsTime().Format(time.RFC822),
		})
	}
	return result
}

func buildAllSeatingAreasResponse(areas []*settingsProto.SeatingArea) []openapi.SeatingArea {
	result := make([]openapi.SeatingArea, 0)
	for _, area := range areas {
		result = append(result, buildSeatingAreaResponse(area))
	}

	return result
}

func buildTableRestrictionsResponse(restrictions []*settingsProto.Restriction) []openapi.TableRestriction {
	result := make([]openapi.TableRestriction, 0)
	for _, r := range restrictions {
		result = append(result, openapi.TableRestriction{
			Date:   r.Date,
			Status: r.Status,
		})
	}

	return result
}

func buildTableResponse(table *settingsProto.Table) openapi.Table {
	return openapi.Table{
		Id:             int32(table.GetId()),
		SeatingAreaId:  int32(table.GetSeatingAreaId()),
		SeatingArea:    buildSeatingAreaResponse(table.GetSeatingArea()),
		TableNumber:    table.GetTableNumber(),
		PosNumber:      int32(table.GetPosNumber()),
		MinPartySize:   int32(table.GetMinPartySize()),
		MaxPartySize:   int32(table.GetMaxPartySize()),
		CombinedTables: table.GetCombinedTables(),
		Restrictions:   buildTableRestrictionsResponse(table.GetRestrictions()),
		CreatedAt:      table.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt:      table.GetUpdatedAt().AsTime().Format(time.RFC822),
	}
}

func buildAllTablesResponse(tables []*settingsProto.Table) []openapi.Table {
	result := make([]openapi.Table, 0)
	for _, table := range tables {
		result = append(result, buildTableResponse(table))
	}

	return result
}

// Country
func buildCountryResponse(country *settingsProto.Country) openapi.Country {
	return openapi.Country{
		Id:                country.GetId(),
		CountryName:       country.GetCountryName(),
		CountryNameArabic: country.GetCountryNameArabic(),
		CountryCode:       country.GetCountryCode(),
		ContinentName:     country.GetContinentName(),
		CountryPhoneCode:  country.CountryPhoneCode,
		TimeZone:          country.GetTimeZone(),
	}
}

func buildAllCountriesResponse(countries []*settingsProto.Country) []openapi.Country {
	var result []openapi.Country
	for _, country := range countries {
		result = append(result, buildCountryResponse(country))
	}

	return result
}

func buildIntegrationResponse(integration *settingsProto.Integration) openapi.Integration {
	return openapi.Integration{
		Id:             integration.GetId(),
		SystemName:     integration.GetSystemName(),
		SystemType:     integration.GetSystemType(),
		BaseURL:        integration.GetBaseURL(),
		CredentialType: integration.GetCredentialType(),
		CreatedAt:      integration.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt:      integration.GetUpdatedAt().AsTime().Format(time.RFC822),
	}
}

func buildAllAutomatedReportTypesResponse(types []*settingsProto.AutomatedReportType) []openapi.AutomatedReportType {
	response := []openapi.AutomatedReportType{}
	for _, reportType := range types {
		response = append(response, openapi.AutomatedReportType{
			Id:          reportType.GetId(),
			Name:        reportType.GetName(),
			Description: reportType.GetDescription(),
		})
	}
	return response
}

func buildCreatorResponse(creator *settingsProto.SCreator) openapi.Creator {
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

func buildAllUsersResponse(users []*settingsProto.SCreator) []openapi.Creator {
	response := []openapi.Creator{}
	for _, user := range users {
		response = append(response, buildCreatorResponse(user))
	}
	return response
}

func buildAutomatedReportResponse(report *settingsProto.AutomatedReport) openapi.AutomatedReport {
	var (
		startTimeStr string
		endTimeStr   string
	)
	if report.GetStartTime() != "" {
		const layout = time.TimeOnly
		parsedTime, err := time.Parse(layout, report.GetStartTime())
		if err != nil {
			log.Printf("Error parsing StartTime: %v", err)
		}

		startTimeWithoutSeconds := time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(),
			parsedTime.Hour(), parsedTime.Minute(), 0, 0, parsedTime.Location())

		const layoutWithoutSeconds = "15:04"
		startTimeStr = startTimeWithoutSeconds.Format(layoutWithoutSeconds)
	}

	if report.GetEndTime() != "" {
		const layout = time.TimeOnly
		parsedTime, err := time.Parse(layout, report.GetEndTime())
		if err != nil {
			log.Printf("Error parsing StartTime: %v", err)
		}

		endTimeWithoutSeconds := time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(),
			parsedTime.Hour(), parsedTime.Minute(), 0, 0, parsedTime.Location())

		const layoutWithoutSeconds = "15:04"
		endTimeStr = endTimeWithoutSeconds.Format(layoutWithoutSeconds)
	}

	res := openapi.AutomatedReport{
		Id:                  report.GetId(),
		Name:                report.GetName(),
		Reports:             buildAllAutomatedReportTypesResponse(report.GetTypes()),
		RepeatInterval:      report.GetRepeatInterval(),
		StartTime:           startTimeStr,
		EndTime:             endTimeStr,
		RepeatIntervalHours: report.GetRepeatIntervalHours(),
		Users:               buildAllUsersResponse(report.GetUsers()),
		SendVia:             strings.Split(report.SendVia, ","),
	}

	return res
}

func buildAllAutomatedReportsResponse(reports []*settingsProto.AutomatedReport) []openapi.AutomatedReport {
	response := []openapi.AutomatedReport{}
	for _, r := range reports {
		response = append(response, buildAutomatedReportResponse(r))
	}
	return response
}

func buildWidgetSettingsResponse(widget *settingsProto.WidgetSettings) openapi.WidgetSettings {
	return openapi.WidgetSettings{
		Font:              widget.GetFont(),
		PrimaryColor:      widget.GetPrimaryColor(),
		SecondaryColor:    widget.GetSecondaryColor(),
		Logo:              widget.GetLogo(),
		Banner:            widget.GetBanner(),
		MainImage:         widget.GetMainImage(),
		FacebookAppId:     widget.GetFacebookAppId(),
		FacebookPixelId:   widget.GetFacebookPixelId(),
		PrivacyPolicyUrl:  widget.GetPrivacyPolicyUrl(),
		TermsOfServiceUrl: widget.GetTermsOfServiceUrl(),
		GoogleAnalyticsId: widget.GetGoogleAnalyticsId(),
		VrUrl:             widget.GetVrUrl(),
	}
}

func buildRestaurantItemResponse(item *settingsProto.RestaurantItem) openapi.RestaurantItem {
	return openapi.RestaurantItem{
		Id:        item.GetId(),
		Name:      item.GetName(),
		Price:     item.GetPrice(),
		Code:      item.GetCode(),
		Image:     item.GetImage(),
		Icon:      item.GetIcon(),
		CreatedAt: item.GetCreatedAt().AsTime().Format(time.RFC3339),
		UpdatedAt: item.GetUpdatedAt().AsTime().Format(time.RFC3339),
	}
}

func buildAllRestaurantItemsResponse(items []*settingsProto.RestaurantItem) []openapi.RestaurantItem {
	response := []openapi.RestaurantItem{}
	for _, item := range items {
		response = append(response, buildRestaurantItemResponse(item))
	}
	return response
}
