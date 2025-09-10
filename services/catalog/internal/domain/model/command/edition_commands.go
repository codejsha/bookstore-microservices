package command

type EditionCreateCommand struct {
	Title          string
	Isbn10         *string
	Isbn13         *string
	NumberOfPages  *int32
	PublishDate    *string
	CoverUids      []string
	Languages      []string
	PhysicalFormat *string
	Description    *string
	WorkUid        string
	PublisherUid   *string
	OlKey          *string
}

type EditionUpdateCommand struct {
	Title          *string
	Isbn10         *string
	Isbn13         *string
	NumberOfPages  *int32
	PublishDate    *string
	CoverUids      []string
	Languages      []string
	PhysicalFormat *string
	Description    *string
	WorkUid        *string
	PublisherUid   *string
	OlKey          *string
}
