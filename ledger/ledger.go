package ledger

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type Ledger struct {
	PrivateKeyFile  string
	WCTLPath        string
	WaveletHost     string
	WaveletPort     uint
	SmartContractID string
}

type Account struct {
	Balance            uint64 `json:"balance"`
	Role               string `json:"role"`
	ReputationReceived []struct {
		Effect   string `json:"effect"`
		ReportID string `json:"report_id"`
	} `json:"reputation_received"`
}

func (l *Ledger) GetReps(accounts []string) (int, error) {
	var total int
	for _, acct := range accounts {
		acctInfo, err := l.callLedgerGetAccount(acct)
		if err != nil {
			return 0, err
		}
		repVal := 0
		for _, rep := range acctInfo.ReputationReceived {
			if rep.Effect == "Positive" {
				repVal++
			}
			if rep.Effect == "Negative" {
				repVal--
			}
		}
		total += repVal
	}
	return total, nil
}

func (l *Ledger) callLedgerGetAccount(accountID string) (*Account, error) {
	payload := fmt.Sprintf(`{
		"account_id": "%s"
	}`, accountID)

	outBytes, err := exec.Command(l.WCTLPath,
		"execute_contract",
		"--api.host",
		l.WaveletHost,
		"--api.port",
		fmt.Sprintf("%d", l.WaveletPort),
		"--api.private_key_file",
		l.PrivateKeyFile,
		l.SmartContractID,
		"fetch_account_info",
		payload).Output()
	if err != nil {
		return nil, err
	}
	var acct Account
	if err := json.Unmarshal(outBytes, &acct); err != nil {
		return nil, err
	}
	return &acct, nil
}
