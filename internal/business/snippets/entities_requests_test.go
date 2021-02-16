package snippets

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_createSnippetRequest_Validate(t *testing.T) {
	now := time.Now().UTC()
	monthAfter := now.Add(time.Hour * 24 * 30)
	hourBefore := now.Add(-time.Hour)
	twoYearsAfter := now.Add(time.Hour * 24 * 365 * 2)

	tests := []struct {
		name      string
		Title     string
		Content   string
		ExpiresAt string
		wantErr   bool
	}{
		{
			name:      "Valid createSnippetRequest",
			Title:     "Valid title",
			Content:   "I want to break free!",
			ExpiresAt: monthAfter.Format(time.RFC3339),
			wantErr:   false,
		},
		{
			name:      "Expires At is too big",
			Title:     "Valid title",
			Content:   "I want to break free!",
			ExpiresAt: twoYearsAfter.Format(time.RFC3339),
			wantErr:   true,
		},
		{
			name:      "Expires At is too small",
			Title:     "Valid title",
			Content:   "I want to break free!",
			ExpiresAt: hourBefore.Format(time.RFC3339),
			wantErr:   true,
		},
		{
			name:      "Empty title",
			Title:     "",
			Content:   "I want to break free!",
			ExpiresAt: monthAfter.Format(time.RFC3339),
			wantErr:   true,
		},
		{
			name:      "Empty content",
			Title:     "Valid title",
			Content:   "",
			ExpiresAt: monthAfter.Format(time.RFC3339),
			wantErr:   true,
		},
		{
			name:      "Too long title",
			Title:     generateLongString(101),
			Content:   "Valid content",
			ExpiresAt: monthAfter.Format(time.RFC3339),
			wantErr:   true,
		},
		{
			name:      "Too long content",
			Title:     "Valid title",
			Content:   generateLongString(10001),
			ExpiresAt: monthAfter.Format(time.RFC3339),
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &createSnippetRequest{
				Title:     tt.Title,
				Content:   tt.Content,
				ExpiresAt: tt.ExpiresAt,
			}
			err := r.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func generateLongString(length int) string {
	buf := bytes.Buffer{}
	for i := 0; i < length; i++ {
		buf.WriteRune('a')
	}
	return buf.String()
}
