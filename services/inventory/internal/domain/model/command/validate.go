package command

import (
	"errors"
	"fmt"
	"strings"
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
