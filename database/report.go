package database

type Report struct {
	ID             string `json:"id"`
	ScammerAddress string `json:"scammer_address"`
	VictimAddress  string `json:"victim_address"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	Proof          string `json:"proof"`
	Timestamp      int64  `json:"timestamp"`
	AccountID      string `json:"account_id"`
	Taint          int    `json:"taint"`
}
