package api

type SubmitReportRequest struct {
	ScammerAddress string `json:"scammer_address" validate:"required"`
	VictimAddress  string `json:"victim_address" validate:"required"`
	Title          string `json:"title" validate:"required"`
	Content        string `json:"content" validate:"required"`
	Proof          string `json:"proof" validate:"required"`
}

type SubmitReportResponse struct {
	ID int `json:"id"`
}
