package profile

import (
	"encoding/json"

	"github.com/joaosoft/web"

	"time"
)

type ErrorResponse struct {
	Code    web.Status `json:"code,omitempty"`
	Message string     `json:"message,omitempty"`
	Cause   string     `json:"cause,omitempty"`
}

type GetSectionRequest struct {
	SectionKey string `json:"section_key" validate:"notzero"`
}

type GetSectionContentsRequest struct {
	SectionKey string `json:"section_key" validate:"notzero"`
}

type Sections []*Section

type Section struct {
	IdSection   string    `json:"id_section" db:"id_section"`
	Key         string    `json:"key" db:"key"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Contents []*Content

type Content struct {
	Key       string           `json:"key" db:"key"`
	Content   *json.RawMessage `json:"content" db:"content"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt time.Time        `json:"updated_at" db:"updated_at"`
}
