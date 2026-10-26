package constant

type TicketStatus string

const (
	TICKET_STATUS_OPEN        TicketStatus = "OPEN"
	TICKET_STATUS_IN_PROGRESS TicketStatus = "IN_PROGRESS"
	TICKET_STATUS_RESOLVED    TicketStatus = "RESOLVED"
	TICKET_STATUS_CLOSED      TicketStatus = "CLOSED"
)

func (s TicketStatus) IsValid() bool {
	switch s {
	case TICKET_STATUS_OPEN, TICKET_STATUS_IN_PROGRESS, TICKET_STATUS_RESOLVED, TICKET_STATUS_CLOSED:
		return true
	}
	return false
}

type TicketPriority string

const (
	TICKET_PRIORITY_LOW    TicketPriority = "LOW"
	TICKET_PRIORITY_MEDIUM TicketPriority = "MEDIUM"
	TICKET_PRIORITY_HIGH   TicketPriority = "HIGH"
	TICKET_PRIORITY_URGENT TicketPriority = "URGENT"
)

func (p TicketPriority) IsValid() bool {
	switch p {
	case TICKET_PRIORITY_LOW, TICKET_PRIORITY_MEDIUM, TICKET_PRIORITY_HIGH, TICKET_PRIORITY_URGENT:
		return true
	}
	return false
}

type AuthorRole string

const (
	AUTHOR_ROLE_CUSTOMER AuthorRole = "CUSTOMER"
	AUTHOR_ROLE_AGENT    AuthorRole = "AGENT"
)

func (r AuthorRole) IsValid() bool {
	switch r {
	case AUTHOR_ROLE_CUSTOMER, AUTHOR_ROLE_AGENT:
		return true
	}
	return false
}
