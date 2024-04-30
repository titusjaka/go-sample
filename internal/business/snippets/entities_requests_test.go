package snippets_test

import (
	"strings"
	"testing"
	"time"

	"github.com/titusjaka/go-sample/internal/business/snippets"
	"github.com/titusjaka/go-sample/internal/infrastructure/utils/testutils"
)

func TestCreateSnippetRequest_Validate(t *testing.T) {
	t.Parallel()

	// Define time variables
	now := time.Now().UTC()
	monthAfter := now.Add(time.Hour * 24 * 30)
	hourBefore := now.Add(-time.Hour)
	twoYearsAfter := now.Add(time.Hour * 24 * 365 * 2)

	tests := []struct {
		name    string
		request snippets.CreateSnippetRequest
		wantErr string
	}{
		{
			name: "Valid CreateSnippetRequest",
			request: snippets.CreateSnippetRequest{
				Title:     "Valid title",
				Content:   "I want to break free!",
				ExpiresAt: monthAfter,
			},
			wantErr: "",
		},
		{
			name: "Invalid: expires_at is too big",
			request: snippets.CreateSnippetRequest{
				Title:     "Valid title",
				Content:   "I want to break free!",
				ExpiresAt: twoYearsAfter,
			},
			wantErr: "expires_at: must be a valid RFC3339 date <= now + 1 year.",
		},
		{
			name: "Invalid: expires_at is too small",
			request: snippets.CreateSnippetRequest{
				Title:     "Valid title",
				Content:   "I want to break free!",
				ExpiresAt: hourBefore,
			},
			wantErr: "expires_at: must be a valid RFC3339 date >= now.",
		},
		{
			name: "Invalid: empty title",
			request: snippets.CreateSnippetRequest{
				Title:     "",
				Content:   "I want to break free!",
				ExpiresAt: monthAfter,
			},
			wantErr: "title: cannot be blank.",
		},
		{
			name: "Invalid: empty content",
			request: snippets.CreateSnippetRequest{
				Title:     "Valid title",
				Content:   "",
				ExpiresAt: monthAfter,
			},
			wantErr: "content: cannot be blank.",
		},
		{
			name: "Invalid: title is too long",
			request: snippets.CreateSnippetRequest{
				Title:     strings.Repeat("a", 101),
				Content:   "Valid content",
				ExpiresAt: monthAfter,
			},
			wantErr: "title: the length must be between 1 and 100.",
		},
		{
			name: "Invalid: content is too long",
			request: snippets.CreateSnippetRequest{
				Title:     "Valid title",
				Content:   strings.Repeat("a", 10001),
				ExpiresAt: monthAfter,
			},
			wantErr: "content: the length must be between 1 and 10000.",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.request.Validate()
			testutils.AssertError(t, tt.wantErr, err)
		})
	}
}
