package auth

import "github.com/M1steryO/Room-Booking-Service/internal/domain"

func (u *AuthUsecase) DummyLogin(role domain.Role) (string, error) {
	if !role.Valid() {
		return "", domain.InvalidRequest("role must be admin or user")
	}

	userID := dummyUserID
	if role == domain.RoleAdmin {
		userID = dummyAdminID
	}

	return u.tokens.Issue(userID, role, u.clock.Now())
}
