package generic

import "github.com/iv-menshenin/accountant/utils/uuid"

type (
	AuthData struct {
		JWT     string    `json:"jwt_token"`
		UserID  uuid.UUID `json:"user_id"`
		Context []string  `json:"context"`
		Refresh string    `json:"refresh_token"`
	}
	AuthQuery struct {
		Login    string `json:"login,omitempty"`
		Password string `json:"password,omitempty"`
	}
	RefreshTokenQuery struct {
		Token string `json:"token,omitempty"`
	}

	Unauthorized struct{}
	Forbidden    struct{}
	NotFound     struct{}
)

func (u Unauthorized) Error() string {
	return "Authentication required"
}

func (f Forbidden) Error() string {
	return "You are not allowed this action"
}

func (n NotFound) Error() string {
	return "Object not found"
}
