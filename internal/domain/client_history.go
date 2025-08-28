package domain

type ClientHistory struct {
	ClientID         string
	OrganizationName string
	HasActiveLoan    bool
	ActiveLoanNumber string
	LastPdn          string
	HasLoans         bool
	ClientFullName   string
}
