package command

const reviewTitleMaxLen = 255

type ReviewWriteCommand struct {
	UserUid string
	BookUid string
	Rating  int32
	Title   *string
	Content *string
}

type ReviewEditCommand struct {
	Rating  *int32
	Title   *string
	Content *string
}

func (c ReviewWriteCommand) Validate() error {
	if err := requireUUID("user_uid", c.UserUid); err != nil {
		return err
	}
	if err := requireUUID("book_uid", c.BookUid); err != nil {
		return err
	}
	if c.Rating < 1 || c.Rating > 5 {
		return invalidf("rating must be between 1 and 5: got %d", c.Rating)
	}
	if err := optionalNonBlank("title", c.Title); err != nil {
		return err
	}
	if err := optionalMaxLen("title", c.Title, reviewTitleMaxLen); err != nil {
		return err
	}
	return optionalNonBlank("content", c.Content)
}

func (c ReviewEditCommand) Validate() error {
	if c.Rating != nil && (*c.Rating < 1 || *c.Rating > 5) {
		return invalidf("rating must be between 1 and 5: got %d", *c.Rating)
	}
	if err := optionalNonBlank("title", c.Title); err != nil {
		return err
	}
	if err := optionalMaxLen("title", c.Title, reviewTitleMaxLen); err != nil {
		return err
	}
	return optionalNonBlank("content", c.Content)
}
