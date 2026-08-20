package aggregate

type PaymentAggregate struct {
	Uid           string
	PaymentUid    string
	CustomerUid   string
	Currency      string
	Amount        float64
	Status        string
	PaymentMethod *string
	ErrorCode     *string
	ErrorMessage  *string
}
