package command

const (
	pointReasonMaxLen   = 255
	pointAmountMaxPerOp = 1_000_000
)

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

func (c PointEarnCommand) Validate() error {
	if err := requireUUID("user_uid", c.UserUid); err != nil {
		return err
	}
	if c.Amount <= 0 {
		return invalidf("amount must be positive: got %d", c.Amount)
	}
	if c.Amount > pointAmountMaxPerOp {
		return invalidf("amount must be at most %d: got %d", pointAmountMaxPerOp, c.Amount)
	}
	if err := optionalNonBlank("reason", c.Reason); err != nil {
		return err
	}
	return optionalMaxLen("reason", c.Reason, pointReasonMaxLen)
}

func (c PointSpendCommand) Validate() error {
	if err := requireUUID("user_uid", c.UserUid); err != nil {
		return err
	}
	if c.Amount <= 0 {
		return invalidf("amount must be positive: got %d", c.Amount)
	}
	if c.Amount > pointAmountMaxPerOp {
		return invalidf("amount must be at most %d: got %d", pointAmountMaxPerOp, c.Amount)
	}
	if err := optionalNonBlank("reason", c.Reason); err != nil {
		return err
	}
	return optionalMaxLen("reason", c.Reason, pointReasonMaxLen)
}
