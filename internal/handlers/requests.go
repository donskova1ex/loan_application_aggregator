package handlers

type CreateOrganizationRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateOrganizationRequest struct {
	Name string `json:"name" validate:"required"`
}

type CreateLoanApplicationRequest struct {
	IncomingOrganizationName string `json:"incoming_organization_name" validate:"required"`
	IssueOrganizationName    string `json:"issue_organization_name" validate:"required"`
	Value                    int64  `json:"value" validate:"required"`
	Phone                    string `json:"phone" validate:"required"`
	Comment                  string `json:"comment"`
}

type UpdateLoanApplicationRequest struct {
	IncomingOrganizationName string `json:"incoming_organization_name"`
	IssueOrganizationName    string `json:"issue_organization_name"`
	Value                    int64  `json:"value"`
	Phone                    string `json:"phone"`
	Comment                  string `json:"comment"`
}
