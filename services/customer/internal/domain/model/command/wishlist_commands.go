package command

import "fmt"

type WishlistAddCommand struct {
	UserUid  string
	BookUids []string
}

type WishlistRemoveCommand struct {
	UserUid  string
	BookUids []string
}

func (c WishlistAddCommand) Validate() error {
	return validateWishlistBooks(c.UserUid, c.BookUids)
}

func (c WishlistRemoveCommand) Validate() error {
	return validateWishlistBooks(c.UserUid, c.BookUids)
}

func validateWishlistBooks(userUid string, bookUids []string) error {
	if err := requireUUID("user_uid", userUid); err != nil {
		return err
	}
	if len(bookUids) == 0 {
		return invalidf("book_uids must not be empty")
	}
	seen := make(map[string]struct{}, len(bookUids))
	for i, bookUid := range bookUids {
		field := fmt.Sprintf("book_uids[%d]", i)
		if err := requireNonBlank(field, bookUid); err != nil {
			return err
		}
		if err := requireUUID(field, bookUid); err != nil {
			return err
		}
		if _, dup := seen[bookUid]; dup {
			return invalidf("book_uids must not contain duplicates: %s", bookUid)
		}
		seen[bookUid] = struct{}{}
	}
	return nil
}
