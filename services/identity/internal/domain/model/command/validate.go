package command

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/codejsha/bookstore-microservices/identity/internal/domain/constant"
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

func requireEmail(field, value string, max int) error {
	if !emailPattern.MatchString(value) {
		return invalidf("%s must be a valid email address", field)
	}
	return requireMaxLen(field, value, max)
}

func requireKnownRoles(field string, roles []string) error {
	seen := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		key := strings.ToUpper(strings.TrimSpace(r))
		if constant.AuthRoleFromString(key) == constant.AUTHROLE_UNKNOWN {
			return invalidf("%s contains unknown role: %q", field, r)
		}
		if _, ok := seen[key]; ok {
			return invalidf("%s must not contain duplicates: %q", field, r)
		}
		seen[key] = struct{}{}
	}
	return nil
}
