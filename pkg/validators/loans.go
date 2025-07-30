package validators

func LoanStatusDescription(statusId int) string {
	switch {
	case statusId == 1:
		return "Специалист"
	case statusId == 2:
		return "ОКР"
	case statusId == 3:
		return "Юрист"
	default:
		return "По умолчанию"
	}
}
