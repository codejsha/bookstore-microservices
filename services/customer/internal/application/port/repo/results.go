package repo

import "time"

type PointResult struct {
	Id      int64
	Uid     string
	UserId  int64
	UserUid string
	Balance int32
}

type PointHistoryResult struct {
	Id         int64
	Uid        string
	UserId     int64
	UserUid    string
	ChangeType string
	Amount     int32
	Reason     *string
	CreatedAt  time.Time
}

type ReviewCreate struct {
	UserUid string
	BookUid string
	Rating  int32
	Title   *string
	Content *string
}

type ReviewResult struct {
	Id        int64
	Uid       string
	UserId    int64
	UserUid   string
	BookId    int64
	BookUid   string
	Rating    int32
	Title     *string
	Content   *string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type WishlistResult struct {
	Id        int64
	Uid       string
	UserId    int64
	UserUid   string
	BookId    int64
	BookUid   string
	CreatedAt time.Time
}
