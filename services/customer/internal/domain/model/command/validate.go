package command

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"

	"github.com/codejsha/bookstore-microservices/customer/internal/domain/constant"
)

var ErrInvalidCommand = errors.New("invalid command")

var emailPattern = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

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
	if len([]rune(value)) > max {
		return invalidf("%s must be at most %d characters", field, max)
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

func optionalEmail(field string, value *string) error {
	if value == nil {
		return nil
	}
	if !emailPattern.MatchString(*value) {
		return invalidf("%s must be a valid email address", field)
	}
	return nil
}

func requireKnownRoles(field string, roles []constant.AuthRole) error {
	seen := make(map[constant.AuthRole]struct{}, len(roles))
	for _, r := range roles {
		if r <= constant.AUTHROLE_UNKNOWN || r > constant.AUTHROLE_VIEW {
			return invalidf("%s contains unknown role: %d", field, r)
		}
		if _, ok := seen[r]; ok {
			return invalidf("%s must not contain duplicates: %d", field, r)
		}
		seen[r] = struct{}{}
	}
	return nil
}
