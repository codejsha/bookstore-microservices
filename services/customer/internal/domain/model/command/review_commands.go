package command

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
