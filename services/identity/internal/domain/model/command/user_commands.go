package command

type UserRegisterCommand struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     *string
	Roles     []string
}

type UserUpdateCommand struct {
	FirstName *string
	LastName  *string
	Phone     *string
}

type UserRolesCommand struct {
	Roles []string
}
