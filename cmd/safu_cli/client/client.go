package client

import (
	"fmt"
	"github.com/perlin-network/noise/crypto"
	"github.com/perlin-network/noise/crypto/ed25519"
	"github.com/pkg/errors"
	"io/ioutil"
	"os/exec"
	"strings"
)

var (
	signaturePolicy = ed25519.New()
)

// Config represents a Perlin Ledger client config.
type Config struct {
	PrivateKeyFile       string
	TaintHost            string
	TaintPort            uint
	WCTLPath             string
	WaveletHost          string
	WaveletPort          uint
	AccountID            string
	SmartContractAddress string
}

// Client represents a Perlin Ledger client.
type Client struct {
	config *Config
}

func NewClient(config *Config) (*Client, error) {
	client := &Client{
		config: config,
	}
	return client, nil
}

func (client *Client) Register() (interface{}, error) {
	return client.callWaveletLedger("custom", fmt.Sprintf(`{
		"receipient": "%s",
		"body": {
			"Payload": "Register"
		}
	}`, client.config.SmartContractAddress))
}

func (client *Client) ResetRep(targetAddress string) (interface{}, error) {
	return client.callWaveletLedger("custom", fmt.Sprintf(`{
		"receipient": "%s",
		"body": {
			"PlusRep": {
				"target_address": "%s"
			}
		}
	}`, client.config.SmartContractAddress, targetAddress))
}

func (client *Client) PlusRep(targetAddress string, scamReportID string) (interface{}, error) {
	return client.callWaveletLedger("custom", fmt.Sprintf(`{
		"receipient": "%s",
		"body": {
			"PlusRep": {
				"target_address": "%s",
				"report_id": "%s"
			}
		}
	}`, client.config.SmartContractAddress, targetAddress, scamReportID))
}

func (client *Client) NegRep(targetAddress string, scamReportID string) (interface{}, error) {
	return client.callWaveletLedger("custom", fmt.Sprintf(`{
		"receipient": "%s",
		"body": {
			"NegRep": {
				"target_address": "%s",
				"report_id": "%s"
			}
		}
	}`, client.config.SmartContractAddress, targetAddress, scamReportID))
}

func (client *Client) Upgrade() (interface{}, error) {
	return client.callWaveletLedger("custom", fmt.Sprintf(`{
		"receipient": "%s",
		"body": {
			"Payload": "UpgradeToVIP"
		}
	}`, client.config.SmartContractAddress))
}

func (client *Client) Deposit(depositAmount int64) (interface{}, error) {
	return client.callWaveletLedger("transfer", fmt.Sprintf(`{
		"receipient": "%s",
		"amount": %d
	}`, client.config.SmartContractAddress, depositAmount))
}

func (client *Client) Query(address string) (interface{}, error) {
	return nil, errors.New("Not implemented")
}

func (client *Client) RegisterScamReport(scammerAddress string, victimAddress string, title string, content string) (interface{}, error) {
	var err error
	// 1. push the scam report to taint server
	_, err = client.callTaintServer()
	if err != nil {
		return nil, err
	}
	scamReportID := "TODO: get this from the taint server"
	// 2. push the report id to the ledger
	_, err = client.callWaveletLedger("custom", fmt.Sprintf(`{
		"receipient": "%s",
		"body": {
			"RegisterScamReport": {
				"report_id": "%s"
			}
		}
	}`, client.config.SmartContractAddress, scamReportID))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"status":    "ok",
		"report_id": scamReportID,
		"tx_id":     "TODO",
	}, nil
}

func getKeyPair(privateKeyFile string) (*crypto.KeyPair, error) {

	bytes, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to open server private key file: %s", privateKeyFile)
	}
	privateKeyHex := strings.TrimSpace(string(bytes))

	return crypto.FromPrivateKey(signaturePolicy, privateKeyHex)
}

func (client *Client) callWaveletLedger(tag string, payload string) (interface{}, error) {
	cmd := fmt.Sprintf("%s send_transaction --api.host %s --api.port %d --api.private_key_file %s %s %s",
		client.config.WCTLPath,
		client.config.WaveletHost,
		client.config.WaveletPort,
		client.config.PrivateKeyFile,
		tag,
		payload,
	)
	_, err := exec.Command(cmd).Output()
	if err != nil {
		return nil, err
	}
	return "TODO", nil
}

func (client *Client) callTaintServer() (interface{}, error) {
	// TODO: take an http request
	return nil, nil
}
