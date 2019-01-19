package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/perlin-network/noise/crypto"
	"github.com/perlin-network/noise/crypto/blake2b"
	"github.com/perlin-network/noise/crypto/ed25519"
	"github.com/perlin-network/safu-go/api"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

var (
	signaturePolicy = ed25519.New()
	hashPolicy      = blake2b.New()
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
		"recipient": "%s",
		"body": {
			"Payload": "Register"
		}
	}`, client.config.SmartContractAddress))
}

func (client *Client) ResetRep(targetAddress string) (interface{}, error) {
	return client.callWaveletLedger("custom", fmt.Sprintf(`{
		"recipient": "%s",
		"body": {
			"PlusRep": {
				"target_address": "%s"
			}
		}
	}`, client.config.SmartContractAddress, targetAddress))
}

func (client *Client) PlusRep(targetAddress string, scamReportID string) (interface{}, error) {
	return client.callWaveletLedger("custom", fmt.Sprintf(`{
		"recipient": "%s",
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
		"recipient": "%s",
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
		"recipient": "%s",
		"body": {
			"Payload": "UpgradeToVIP"
		}
	}`, client.config.SmartContractAddress))
}

func (client *Client) Deposit(depositAmount int64) (interface{}, error) {
	return client.callWaveletLedger("transfer", fmt.Sprintf(`{
		"recipient": "%s",
		"amount": %d
	}`, client.config.SmartContractAddress, depositAmount))
}

func (client *Client) Query(targetAddress string) (interface{}, error) {
	body := api.QueryAddressRequest{
		AccountID:     client.config.AccountID,
		Timestamp:     time.Now().UnixNano(),
		TargetAddress: targetAddress,
	}
	var taintResp api.QueryAddressResponse
	if err := client.callTaintServer(api.RouteQueryAddress, body, &taintResp); err != nil {
		return nil, err
	}
	return taintResp, nil
}

func (client *Client) RegisterScamReport(scammerAddress string, victimAddress string, title string, content string) (interface{}, error) {
	var err error
	// 1. push the scam report to taint server
	body := api.SubmitReportRequest{
		AccountID:      client.config.AccountID,
		Timestamp:      time.Now().UnixNano(),
		ScammerAddress: scammerAddress,
		VictimAddress:  victimAddress,
		Title:          title,
		Content:        content,
	}

	keyPair, err := getKeyPair(client.config.PrivateKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get private key")
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to marshal body")
	}

	sig, err := keyPair.Sign(signaturePolicy, hashPolicy, bodyBytes)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to sign body")
	}

	body.Proof = base64.StdEncoding.EncodeToString(sig)

	var taintResp api.SubmitReportResponse
	if err = client.callTaintServer(api.RoutePostScamRepot, body, &taintResp); err != nil {
		return nil, err
	}

	// 2. push the report id to the ledger
	ledgerResp, err := client.callWaveletLedger("custom", fmt.Sprintf(`{
		"recipient": "%s",
		"body": {
			"RegisterScamReport": {
				"report_id": "%s"
			}
		}
	}`, client.config.SmartContractAddress, taintResp.ID))
	if err != nil {
		return nil, err
	}

	ledgerUnmashalled, ok := ledgerResp.(map[string]interface{})
	if !ok {
		return nil, errors.New("Unable to cast transaction")
	}
	txID, ok := ledgerUnmashalled["transaction_id"].(string)
	if !ok {
		return nil, errors.New("Unable to parse transaction id")
	}

	return map[string]string{
		"status":    "ok",
		"report_id": taintResp.ID,
		"tx_id":     txID,
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
	outBytes, err := exec.Command(client.config.WCTLPath,
		"send_transaction",
		"--api.host",
		client.config.WaveletHost,
		"--api.port",
		fmt.Sprintf("%d", client.config.WaveletPort),
		"--api.private_key_file",
		client.config.PrivateKeyFile,
		tag,
		payload).Output()
	if err != nil {
		return nil, err
	}
	var result interface{}
	json.Unmarshal(outBytes, &result)
	return result, nil
}

func (client *Client) callTaintServer(path string, body interface{}, out interface{}) error {
	prot := "http"
	u, err := url.Parse(fmt.Sprintf("%s://%s:%d%s", prot, client.config.TaintHost, client.config.TaintPort, path))
	if err != nil {
		return err
	}
	req := &http.Request{
		Method: "POST",
		URL:    u,
	}

	if body != nil {
		rawBody, err := json.Marshal(body)
		if err != nil {
			return err
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(rawBody))
	}

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.Errorf("got an error code %v: %v", resp.Status, string(data))
	}

	if out == nil {
		return nil
	}
	return json.Unmarshal(data, out)
}
