package command

type SubjectCreateCommand struct {
	Name string
}

type SubjectUpdateCommand struct {
	Name *string
}
