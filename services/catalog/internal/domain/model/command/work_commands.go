package command

type WorkCreateCommand struct {
	Title            string
	Description      *string
	CoverUids        []string
	FirstPublishDate *string
	OlKey            *string
	AuthorUids       []string
	SubjectNames     []string
}

type WorkUpdateCommand struct {
	Title            *string
	Description      *string
	CoverUids        []string
	FirstPublishDate *string
	OlKey            *string
	AuthorUids       []string
	SubjectNames     []string
}
