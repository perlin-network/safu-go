package api

type SubmitReportRequest struct {
	Timestamp      int64  `json:"timestamp" validate:"required"`
	AccountID      string `json:"account_id" validate:"required"`
	ScammerAddress string `json:"scammer_address" validate:"required"`
	VictimAddress  string `json:"victim_address" validate:"required"`
	Title          string `json:"title" validate:"required"`
	Content        string `json:"content" validate:"required"`
	Proof          string `json:"proof"`
}

type SubmitReportResponse struct {
	ID string `json:"id"`
}

type QueryAddressRequest struct {
	Timestamp     int64  `json:"timestamp" validate:"required"`
	AccountID     string `json:"account_id" validate:"required"`
	TargetAddress string `json:"target_address" validate:"required"`
}

type QueryAddressResponse struct {
	TargetAddress string `json:"target_address" validate:"required"`
	TaintScore    int32  `json:"taint_score" validate:"required"`
}

type ScamReport struct {
	ID             string `json:"id"`
	Timestamp      int64  `json:"timestamp" validate:"required"`
	AccountID      string `json:"account_id" validate:"required"`
	ScammerAddress string `json:"scammer_address" validate:"required"`
	VictimAddress  string `json:"victim_address" validate:"required"`
	Title          string `json:"title" validate:"required"`
	Content        string `json:"content" validate:"required"`
	Proof          string `json:"proof" validate:"required"`
	Taint          int    `json:"taint"`
}

type AllScamReportResponse struct {
	Reports []*ScamReport `json:"reports"`
}
