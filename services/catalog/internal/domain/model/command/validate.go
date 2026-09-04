package command

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

var ErrInvalidCommand = errors.New("invalid command")

func invalidf(format string, args ...any) error {
	return fmt.Errorf(format+": %w", append(args, ErrInvalidCommand)...)
}

func requireNonBlank(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return invalidf("%s must not be blank", field)
	}
	return nil
}

func optionalNonBlank(field string, value *string) error {
	if value == nil {
		return nil
	}
	return requireNonBlank(field, *value)
}

func requireMaxLen(field, value string, max int) error {
	if utf8.RuneCountInString(value) > max {
		return invalidf("%s must not exceed %d characters", field, max)
	}
	return nil
}

func optionalMaxLen(field string, value *string, max int) error {
	if value == nil {
		return nil
	}
	return requireMaxLen(field, *value, max)
}

func requireUUID(field, value string) error {
	if _, err := uuid.Parse(value); err != nil {
		return invalidf("%s must be a valid uuid", field)
	}
	return nil
}

func optionalUUID(field string, value *string) error {
	if value == nil {
		return nil
	}
	return requireUUID(field, *value)
}

func requireElementsUUID(field string, values []string) error {
	for _, v := range values {
		if _, err := uuid.Parse(v); err != nil {
			return invalidf("%s must contain only valid uuids", field)
		}
	}
	return nil
}

func requireElementsNonBlank(field string, values []string) error {
	for _, v := range values {
		if strings.TrimSpace(v) == "" {
			return invalidf("%s must not contain blank entries", field)
		}
	}
	return nil
}

func requireElementsMaxLen(field string, values []string, max int) error {
	for _, v := range values {
		if utf8.RuneCountInString(v) > max {
			return invalidf("%s entries must not exceed %d characters", field, max)
		}
	}
	return nil
}

func requireUniqueElements(field string, values []string) error {
	seen := make(map[string]struct{}, len(values))
	for _, v := range values {
		if _, ok := seen[v]; ok {
			return invalidf("%s must not contain duplicates: %s", field, v)
		}
		seen[v] = struct{}{}
	}
	return nil
}

func optionalPositive(field string, value *int32) error {
	if value != nil && *value <= 0 {
		return invalidf("%s must be positive: got %d", field, *value)
	}
	return nil
}

func validateIsbn10(value *string) error {
	if value == nil {
		return nil
	}
	s := *value
	if len(s) != 10 {
		return invalidf("isbn10 must be 10 characters")
	}
	sum := 0
	for i, r := range s {
		var d int
		switch {
		case r >= '0' && r <= '9':
			d = int(r - '0')
		case (r == 'X' || r == 'x') && i == 9:
			d = 10
		default:
			return invalidf("isbn10 has invalid character %q", r)
		}
		sum += (10 - i) * d
	}
	if sum%11 != 0 {
		return invalidf("isbn10 checksum mismatch")
	}
	return nil
}

func validateIsbn13(value *string) error {
	if value == nil {
		return nil
	}
	s := *value
	if len(s) != 13 {
		return invalidf("isbn13 must be 13 digits")
	}
	if !strings.HasPrefix(s, "978") && !strings.HasPrefix(s, "979") {
		return invalidf("isbn13 must start with 978 or 979")
	}
	sum := 0
	for i, r := range s {
		if r < '0' || r > '9' {
			return invalidf("isbn13 has invalid character %q", r)
		}
		d := int(r - '0')
		if i%2 == 1 {
			d *= 3
		}
		sum += d
	}
	if sum%10 != 0 {
		return invalidf("isbn13 checksum mismatch")
	}
	return nil
}

func validateLifespan(birthDate, deathDate *string) error {
	if birthDate == nil || deathDate == nil {
		return nil
	}
	birth, err := time.Parse(time.DateOnly, *birthDate)
	if err != nil {
		return nil
	}
	death, err := time.Parse(time.DateOnly, *deathDate)
	if err != nil {
		return nil
	}
	if death.Before(birth) {
		return invalidf("death_date must not precede birth_date")
	}
	return nil
}
