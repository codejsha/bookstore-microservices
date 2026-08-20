package command

type AuthorCreateCommand struct {
	Name           string
	Bio            *string
	BirthDate      *string
	DeathDate      *string
	PhotoUids      []string
	AlternateNames []string
	OlKey          *string
}

type AuthorUpdateCommand struct {
	Name           *string
	Bio            *string
	BirthDate      *string
	DeathDate      *string
	PhotoUids      []string
	AlternateNames []string
	OlKey          *string
}
