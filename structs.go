package session

import (
	"github.com/joaosoft/types"
	"github.com/joaosoft/web"

	"time"
)

type ErrorResponse struct {
	Code    web.Status `json:"code,omitempty"`
	Message string     `json:"message,omitempty"`
	Cause   string     `json:"cause,omitempty"`
}

type GetSessionRequest struct {
	Email    string `json:"email" validate:"notzero"`
	Password string `json:"password" validate:"notzero"`
}

type RefreshSessionRequest struct {
	Authorization string `json:"authorization" validate:"notzero"`
}

type SessionResponse struct {
	TokenType    string `json:"token_type"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type UpdateProcessRequest struct {
	IdProcess string `json:"id_process" validate:"notzero"`

	Body struct {
		Type           string         `json:"type" validate:"notzero"`
		Name           string         `json:"name" validate:"notzero"`
		Description    string         `json:"description"`
		DateFrom       *types.Date    `json:"date_from" validate:"special={date}"`
		DateTo         *types.Date    `json:"date_to" validate:"special={date}"`
		TimeFrom       *types.Time    `json:"time_from" validate:"special={time}"`
		TimeTo         *types.Time    `json:"time_to" validate:"special={time}"`
		DaysOff        *types.ListDay `json:"days_off" validate:"options=monday;tuesday;wednesday;thursday;friday;saturday;sunday"`
		Authentication string         `json:"monitor"`
	}
}

type User struct {
	IdUser       string    `json:"id_user" db:"id_user"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	Email        string    `json:"email" db:"email"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	Active       bool      `json:"active" db:"active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
