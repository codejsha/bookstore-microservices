package command

type PublisherCreateCommand struct {
	Name    string
	Address *string
	OlKey   *string
}

type PublisherUpdateCommand struct {
	Name    *string
	Address *string
	OlKey   *string
}
