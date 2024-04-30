package snippets

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ListSnippetsRequest represents a request struct for GET /snippets?limit=<x>&offset=<y> method
type ListSnippetsRequest struct {
	Limit  uint `schema:"limit"`
	Offset uint `schema:"offset"`
}

// CreateSnippetRequest represents a request struct for POST /snippets method
type CreateSnippetRequest struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Validate implements ozzo-validation.Validatable interface and used to check user request
func (r *CreateSnippetRequest) Validate() error {
	now := time.Now().UTC().Truncate(time.Second)

	yearAfter := now.Add(366 * 24 * time.Hour)
	rules := []*validation.FieldRules{
		validation.Field(&r.Title, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Content, validation.Required, validation.Length(1, 10000)),
		validation.Field(
			&r.ExpiresAt,
			validation.Required,
			validation.Min(now).Error("must be a valid RFC3339 date >= now"),
			validation.Max(yearAfter).Error("must be a valid RFC3339 date <= now + 1 year"),
		),
	}

	return validation.ValidateStruct(r, rules...)
}
