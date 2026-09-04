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

func (c EditionCreateCommand) Validate() error {
	if err := requireNonBlank("title", c.Title); err != nil {
		return err
	}
	if err := requireMaxLen("title", c.Title, 512); err != nil {
		return err
	}
	if err := requireNonBlank("work_uid", c.WorkUid); err != nil {
		return err
	}
	if err := requireUUID("work_uid", c.WorkUid); err != nil {
		return err
	}
	if err := optionalUUID("publisher_uid", c.PublisherUid); err != nil {
		return err
	}
	if err := optionalMaxLen("publish_date", c.PublishDate, 32); err != nil {
		return err
	}
	if err := optionalMaxLen("physical_format", c.PhysicalFormat, 64); err != nil {
		return err
	}
	if err := optionalNonBlank("ol_key", c.OlKey); err != nil {
		return err
	}
	if err := optionalMaxLen("ol_key", c.OlKey, 64); err != nil {
		return err
	}
	if err := validateIsbn10(c.Isbn10); err != nil {
		return err
	}
	if err := validateIsbn13(c.Isbn13); err != nil {
		return err
	}
	return optionalPositive("number_of_pages", c.NumberOfPages)
}

func (c EditionUpdateCommand) Validate() error {
	if err := optionalNonBlank("title", c.Title); err != nil {
		return err
	}
	if err := optionalMaxLen("title", c.Title, 512); err != nil {
		return err
	}
	if err := optionalNonBlank("work_uid", c.WorkUid); err != nil {
		return err
	}
	if err := optionalUUID("work_uid", c.WorkUid); err != nil {
		return err
	}
	if err := optionalUUID("publisher_uid", c.PublisherUid); err != nil {
		return err
	}
	if err := optionalMaxLen("publish_date", c.PublishDate, 32); err != nil {
		return err
	}
	if err := optionalMaxLen("physical_format", c.PhysicalFormat, 64); err != nil {
		return err
	}
	if err := optionalNonBlank("ol_key", c.OlKey); err != nil {
		return err
	}
	if err := optionalMaxLen("ol_key", c.OlKey, 64); err != nil {
		return err
	}
	if err := validateIsbn10(c.Isbn10); err != nil {
		return err
	}
	if err := validateIsbn13(c.Isbn13); err != nil {
		return err
	}
	return optionalPositive("number_of_pages", c.NumberOfPages)
}
