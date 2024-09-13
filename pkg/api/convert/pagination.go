package convert

import (
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func BuildPaginationResponse(proto *guestProto.Pagination) openapi.Pagination {
	return openapi.Pagination{
		Total:       proto.GetTotal(),
		PerPage:     proto.GetPerPage(),
		CurrentPage: proto.GetCurrentPage(),
		LastPage:    proto.GetLastPage(),
		From:        proto.GetFrom(),
		To:          proto.GetTo(),
	}
}
