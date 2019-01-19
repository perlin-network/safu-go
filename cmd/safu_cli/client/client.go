package client

import (
	"github.com/perlin-network/noise/crypto"
	"github.com/pkg/errors"
)

// Config represents a Perlin Ledger client config.
type Config struct {
	Host          string
	Port          uint
	PrivateKeyHex string
}

// Client represents a Perlin Ledger client.
type Client struct {
	Config  Config
	KeyPair *crypto.KeyPair
}

func NewClient(config Config) (*Client, error) {
	return nil, nil
}

func (client *Client) Register() (interface{}, error) {
	return nil, errors.New("Not implemented")
}

func (client *Client) ResetRep(targetAddress string) (interface{}, error) {
	return nil, errors.New("Not implemented")
}

func (client *Client) PlusRep(address string, scamReportID string) (interface{}, error) {
	return nil, errors.New("Not implemented")
}

func (client *Client) NegRep(address string, scamReportID string) (interface{}, error) {
	return nil, errors.New("Not implemented")
}

func (client *Client) Upgrade() (interface{}, error) {
	return nil, errors.New("Not implemented")
}

func (client *Client) Deposit(depositAmount int) (interface{}, error) {
	return nil, errors.New("Not implemented")
}

func (client *Client) Withdraw(withdrawAmount int) (interface{}, error) {
	return nil, errors.New("Not implemented")
}

func (client *Client) Balance() (interface{}, error) {
	return nil, errors.New("Not implemented")
}

func (client *Client) Query(address string) (interface{}, error) {
	return nil, errors.New("Not implemented")
}

func (client *Client) RegisterScamReport(scamReportID string) (interface{}, error) {
	return nil, errors.New("Not implemented")
}
