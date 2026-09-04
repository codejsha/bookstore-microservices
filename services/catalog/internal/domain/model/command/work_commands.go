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

func (c WorkCreateCommand) Validate() error {
	if err := requireNonBlank("title", c.Title); err != nil {
		return err
	}
	if err := requireMaxLen("title", c.Title, 512); err != nil {
		return err
	}
	if err := optionalMaxLen("first_publish_date", c.FirstPublishDate, 32); err != nil {
		return err
	}
	if err := optionalNonBlank("ol_key", c.OlKey); err != nil {
		return err
	}
	if err := optionalMaxLen("ol_key", c.OlKey, 64); err != nil {
		return err
	}
	if len(c.AuthorUids) == 0 {
		return invalidf("author_uids must not be empty")
	}
	if err := requireElementsNonBlank("author_uids", c.AuthorUids); err != nil {
		return err
	}
	if err := requireElementsUUID("author_uids", c.AuthorUids); err != nil {
		return err
	}
	if err := requireUniqueElements("author_uids", c.AuthorUids); err != nil {
		return err
	}
	if err := requireElementsNonBlank("subject_names", c.SubjectNames); err != nil {
		return err
	}
	if err := requireElementsMaxLen("subject_names", c.SubjectNames, 255); err != nil {
		return err
	}
	return requireUniqueElements("subject_names", c.SubjectNames)
}

func (c WorkUpdateCommand) Validate() error {
	if err := optionalNonBlank("title", c.Title); err != nil {
		return err
	}
	if err := optionalMaxLen("title", c.Title, 512); err != nil {
		return err
	}
	if err := optionalMaxLen("first_publish_date", c.FirstPublishDate, 32); err != nil {
		return err
	}
	if err := optionalNonBlank("ol_key", c.OlKey); err != nil {
		return err
	}
	if err := optionalMaxLen("ol_key", c.OlKey, 64); err != nil {
		return err
	}
	if err := requireElementsNonBlank("author_uids", c.AuthorUids); err != nil {
		return err
	}
	if err := requireElementsUUID("author_uids", c.AuthorUids); err != nil {
		return err
	}
	if err := requireUniqueElements("author_uids", c.AuthorUids); err != nil {
		return err
	}
	if err := requireElementsNonBlank("subject_names", c.SubjectNames); err != nil {
		return err
	}
	if err := requireElementsMaxLen("subject_names", c.SubjectNames, 255); err != nil {
		return err
	}
	return requireUniqueElements("subject_names", c.SubjectNames)
}
