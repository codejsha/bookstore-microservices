package command

type PointEarnCommand struct {
	UserUid string
	Amount  int32
	Reason  *string
}

type PointSpendCommand struct {
	UserUid string
	Amount  int32
	Reason  *string
}
