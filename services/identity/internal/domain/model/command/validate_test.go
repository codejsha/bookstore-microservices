package command

import (
	"errors"
	"strings"
	"testing"
)

func ptr[T any](v T) *T { return &v }

func TestUserRegisterCommand_Validate(t *testing.T) {
	valid := UserRegisterCommand{Email: "a@b.co", Password: "secret", FirstName: "Jin", LastName: "Ha"}
	tests := []struct {
		name    string
		mutate  func(c *UserRegisterCommand)
		wantErr bool
	}{
		{"whenOnlyRequiredFieldsSet_returnsNil", func(c *UserRegisterCommand) {}, false},
		{"whenRolesKnown_returnsNil", func(c *UserRegisterCommand) { c.Roles = []string{"PROFILE", "VIEW"} }, false},
		{"whenEmailMalformed_returnsErrInvalidCommand", func(c *UserRegisterCommand) { c.Email = "not-an-email" }, true},
		{"whenPasswordBlank_returnsErrInvalidCommand", func(c *UserRegisterCommand) { c.Password = " " }, true},
		{"whenFirstNameBlank_returnsErrInvalidCommand", func(c *UserRegisterCommand) { c.FirstName = "" }, true},
		{"whenLastNameBlank_returnsErrInvalidCommand", func(c *UserRegisterCommand) { c.LastName = "  " }, true},
		{"whenPhoneBlank_returnsErrInvalidCommand", func(c *UserRegisterCommand) { c.Phone = ptr(" ") }, true},
		{"whenRoleUnknown_returnsErrInvalidCommand", func(c *UserRegisterCommand) { c.Roles = []string{"SUPERUSER"} }, true},
		{"whenRoleIsUnknownLiteral_returnsErrInvalidCommand", func(c *UserRegisterCommand) { c.Roles = []string{"UNKNOWN"} }, true},
		{"whenRolesDuplicated_returnsErrInvalidCommand", func(c *UserRegisterCommand) { c.Roles = []string{"VIEW", "view"} }, true},
		{"whenEmailAtLimit_returnsNil", func(c *UserRegisterCommand) { c.Email = strings.Repeat("a", 95) + "@b.co" }, false},
		{"whenEmailOverLimit_returnsErrInvalidCommand", func(c *UserRegisterCommand) { c.Email = strings.Repeat("a", 96) + "@b.co" }, true},
		{"whenFirstNameAtLimit_returnsNil", func(c *UserRegisterCommand) { c.FirstName = strings.Repeat("j", 50) }, false},
		{"whenFirstNameOverLimit_returnsErrInvalidCommand", func(c *UserRegisterCommand) { c.FirstName = strings.Repeat("j", 51) }, true},
		{"whenLastNameOverLimit_returnsErrInvalidCommand", func(c *UserRegisterCommand) { c.LastName = strings.Repeat("h", 51) }, true},
		{"whenPhoneAtLimit_returnsNil", func(c *UserRegisterCommand) { c.Phone = ptr(strings.Repeat("1", 30)) }, false},
		{"whenPhoneOverLimit_returnsErrInvalidCommand", func(c *UserRegisterCommand) { c.Phone = ptr(strings.Repeat("1", 31)) }, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := valid
			tt.mutate(&cmd)
			err := cmd.Validate()
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidCommand) {
					t.Fatalf("want ErrInvalidCommand, got %v", err)
				}
			} else if err != nil {
				t.Fatalf("want nil, got %v", err)
			}
		})
	}
}

// A partial update carries only the fields the caller set: an empty command passes and
// only a present-but-blank or over-limit field is rejected.
func TestUserUpdateCommand_WhenFieldBlankOrOverLimit_ReturnsErrInvalidCommand(t *testing.T) {
	if err := (UserUpdateCommand{}).Validate(); err != nil {
		t.Fatalf("empty update must be valid, got %v", err)
	}
	if err := (UserUpdateCommand{FirstName: ptr("Jin")}).Validate(); err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if err := (UserUpdateCommand{Phone: ptr(" ")}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand, got %v", err)
	}
	if err := (UserUpdateCommand{FirstName: ptr(strings.Repeat("j", 51))}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("first name over limit: want ErrInvalidCommand, got %v", err)
	}
	if err := (UserUpdateCommand{LastName: ptr(strings.Repeat("h", 51))}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("last name over limit: want ErrInvalidCommand, got %v", err)
	}
	if err := (UserUpdateCommand{Phone: ptr(strings.Repeat("1", 31))}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("phone over limit: want ErrInvalidCommand, got %v", err)
	}
	if err := (UserUpdateCommand{FirstName: ptr(strings.Repeat("j", 50)), Phone: ptr(strings.Repeat("1", 30))}).Validate(); err != nil {
		t.Fatalf("at limit: want nil, got %v", err)
	}
}

// One pass over the role rules: a known set passes, an empty set, an unknown role and
// a duplicated role are each rejected.
func TestUserRolesCommand_WhenRolesEmptyUnknownOrDuplicated_ReturnsErrInvalidCommand(t *testing.T) {
	if err := (UserRolesCommand{Roles: []string{"ADMIN", "MANAGE"}}).Validate(); err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if err := (UserRolesCommand{}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("empty roles: want ErrInvalidCommand, got %v", err)
	}
	if err := (UserRolesCommand{Roles: []string{"NOPE"}}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("unknown role: want ErrInvalidCommand, got %v", err)
	}
	if err := (UserRolesCommand{Roles: []string{"VIEW", "VIEW"}}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("duplicate roles: want ErrInvalidCommand, got %v", err)
	}
}
