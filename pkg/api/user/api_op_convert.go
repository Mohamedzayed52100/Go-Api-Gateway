package user

import (
	"time"

	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func buildPinCodeResponse(proto *userProto.PinCodeResult) *openapi.PinCodeResponse {
	return &openapi.PinCodeResponse{
		Status:  proto.GetStatus(),
		Message: proto.GetStatus(),
	}
}

func buildAuthenticatedUserResponse(user *userProto.AuthenticatedUser) openapi.AuthenticatedUser {
	return openapi.AuthenticatedUser{
		Id:          user.GetId(),
		EmployeeId:  user.GetEmployeeId(),
		BirthDate:   user.GetBirthDate().AsTime().Format(time.DateOnly),
		FirstName:   user.GetFirstName(),
		LastName:    user.GetLastName(),
		Email:       user.GetEmail(),
		PhoneNumber: user.GetPhoneNumber(),
		Gender:      user.GetGender(),
		Role:        buildRoleResponse(user.GetRole()),
		Department:  buildDepartmentResponse(user.GetDepartment()),
		Timezone:    user.GetTimezone(),
		Avatar:      user.GetAvatar(),
		Branch: openapi.ReservationBranch{
			Id:       user.GetBranch().GetId(),
			Name:     user.GetBranch().GetName(),
			Currency: user.GetBranch().GetCurrency(),
		},
		GroupId: user.GetGroupId(),
	}
}

func buildUserResponse(user *userProto.User) openapi.User {
	res := openapi.User{
		Id:          user.GetId(),
		EmployeeId:  user.GetEmployeeId(),
		Email:       user.GetEmail(),
		FirstName:   user.GetFirstName(),
		LastName:    user.GetLastName(),
		PhoneNumber: user.GetPhoneNumber(),
		Gender:      user.GetGender(),
		Avatar:      user.GetAvatar(),
		Role:        user.GetRole(),
		Department:  user.GetDepartment(),
		JoinedAt:    user.GetJoinedAt().AsTime().Format(time.DateOnly),
		Birthdate:   user.GetBirthdate().AsTime().Format(time.DateOnly),
	}

	if user.GetBranches() != nil {
		for _, branch := range user.Branches {
			res.Branches = append(res.Branches, openapi.UserBranch{
				Id:   int32(branch.Id),
				Name: branch.Name,
			})
		}
	}

	return res
}

func buildAllUsersResponse(users []*userProto.User) []openapi.User {
	res := []openapi.User{}
	for _, user := range users {
		res = append(res, buildUserResponse(user))
	}
	return res
}

func buildPaginationResponse(pagination *userProto.UPagination) openapi.Pagination {
	return openapi.Pagination{
		Total:       pagination.GetTotal(),
		PerPage:     pagination.GetPerPage(),
		CurrentPage: pagination.GetCurrentPage(),
		LastPage:    pagination.GetLastPage(),
		From:        pagination.GetFrom(),
		To:          pagination.GetTo(),
	}
}

func buildRoleResponse(role *userProto.Role) openapi.Role {
	res := openapi.Role{
		Id:          role.GetId(),
		Name:        role.GetName(),
		DisplayName: role.GetDisplayName(),
		Department:  role.GetDepartment(),
		UsersCount:  role.GetUsersCount(),
		CreatedAt:   role.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt:   role.GetUpdatedAt().AsTime().Format(time.RFC822),
	}

	if role.GetPermissions() != nil {
		res.Permissions = role.GetPermissions()
	} else {
		res.Permissions = []string{}
	}

	return res
}

func buildAllRolesResponse(roles []*userProto.Role) []openapi.Role {
	res := []openapi.Role{}
	for _, role := range roles {
		res = append(res, buildRoleResponse(role))
	}
	return res
}

func buildPermissionResponse(permission *userProto.Permission) openapi.Permission {
	res := openapi.Permission{
		Id:          permission.GetId(),
		Name:        permission.GetName(),
		DisplayName: permission.GetDisplayName(),
		Description: permission.GetDescription(),
		CreatedAt:   permission.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt:   permission.GetUpdatedAt().AsTime().Format(time.RFC822),
	}

	if permission.GetPermissions() != nil {
		conv := buildAllPermissionsResponse(permission.GetPermissions())
		res.Permissions = conv
	} else {
		res.Permissions = nil
	}

	return res
}

func buildAllPermissionsResponse(permissions []*userProto.Permission) []openapi.Permission {
	res := []openapi.Permission{}

	for _, permission := range permissions {
		res = append(res, buildPermissionResponse(permission))
	}

	return res
}

func buildDepartmentResponse(department *userProto.Department) openapi.Department {
	return openapi.Department{
		Id:        department.GetId(),
		Name:      department.GetName(),
		CreatedAt: department.GetCreatedAt().AsTime().Format(time.RFC822),
		UpdatedAt: department.GetUpdatedAt().AsTime().Format(time.RFC822),
	}
}

func buildAllDepartmentsResponse(departments []*userProto.Department) []openapi.Department {
	res := []openapi.Department{}
	for _, department := range departments {
		res = append(res, buildDepartmentResponse(department))
	}
	return res
}

// Branch
func buildBranchResponse(branch *userProto.Branch) openapi.Branch {
	return openapi.Branch{
		Id:          branch.GetId(),
		Name:        branch.GetName(),
		Country:     branch.GetCountry(),
		City:        branch.GetCity(),
		Address:     branch.GetAddress(),
		GmapsLink:   branch.GetGmapsLink(),
		Email:       branch.GetEmail(),
		PhoneNumber: branch.GetPhoneNumber(),
		Website:     branch.GetWebsite(),
	}
}

func buildAllBranchesResponse(branches []*userProto.Branch) []openapi.Branch {
	result := make([]openapi.Branch, 0)
	for _, branch := range branches {
		result = append(result, buildBranchResponse(branch))
	}

	return result
}
