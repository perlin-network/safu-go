package api

type SubmitReportRequest struct {
	ScammerAddress string `json:"scammer_address"`
	VictimAddress  string `json:"victim_address"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	Proof          string `json:"proof"`
}

type SubmitReportResponse struct {
}
