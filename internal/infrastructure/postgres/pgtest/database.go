package pgtest

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"testing"
	"unicode"
)

// DatabaseName creates a new database name with a random suffix.
func DatabaseName(t testing.TB) string {
	t.Helper()

	random := rand.IntN(99999999-10000000) + 10000000 // nolint:gosec

	// Use testdb_<random>_<test name> as the database name.
	// <random> should go first to avoid conflicts with the 63 characters limit.
	name := strings.ToLower(fmt.Sprintf("testdb_%d_%s", random, t.Name()))

	name = strings.Map(func(r rune) rune {
		if !unicode.In(r, unicode.Letter, unicode.Digit) {
			return '_'
		}
		return r
	}, name)

	// PostgreSQL database names are limited to 63 characters.
	name = name[:min(len(name), 63)]

	return name
}
