package snippets

import (
	"time"
)

// Snippet model struct
type Snippet struct {
	ID        uint
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}
