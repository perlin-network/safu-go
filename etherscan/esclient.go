package etherscan

import (
	"github.com/nanmu42/etherscan-api"
	"github.com/perlin-network/safu-go/model"
	"log"
	"strings"
)

// Etherscan client
type ESClient struct {
	client *etherscan.Client
}

func NewESClient(APIKey string) *ESClient {
	return &ESClient{
		client: etherscan.New(etherscan.Mainnet, APIKey),
	}
}

// WIP
func (e *ESClient) Crawl(address string) ([]*model.Vertex, error) {
	var list []*model.Vertex

	txs, err := e.getTxByAddress(address, 1, 10)
	if err != nil {
		return nil, err
	}

	var maxDepth = 3
	var currentDepth = 0
	var elementsToDepthIncrease = 1
	var nextElementsToDepthIncrease = 0

	visited := make(map[string]bool)
	visited[address] = true

	var q []etherscan.NormalTx
	q = append(q, txs...)

	for len(q) != 0 {
		u := q[0]
		q = q[1:len(q):len(q)]

		toTxs, err := e.getTxByAddress(u.To, 1, 10)
		log.Printf("q: %d, depth: %d, get: %s, result:= %d", len(q), currentDepth, u.To, len(toTxs))

		if err != nil {
			continue
		}

		elementsToDepthIncrease--
		if elementsToDepthIncrease == 0 {
			currentDepth++
			if currentDepth > maxDepth {
				log.Printf("reached max depth %d", maxDepth)
				return list, nil
			}

			elementsToDepthIncrease = nextElementsToDepthIncrease
			nextElementsToDepthIncrease = 0
		}

		for _, toTx := range toTxs {
			// Ignore if it's not from the address
			if toTx.From != address {
				log.Println("not from address")
				continue
			}

			// Check if we've already visited
			if _, ok := visited[toTx.To]; ok {
				log.Println("already visited")
				continue
			}

			nextElementsToDepthIncrease++
			visited[toTx.To] = true
			q = append(q, toTx)
		}


	}

	return list, nil
}

// WIP
func (e *ESClient) getTxByAddress(address string, page int, offset int) ([]etherscan.NormalTx, error) {
	startBlock := 0
	endBlock := 99999999

	log.Printf("getTxByAddress address: %s, page: %d, offset: %d", address, page, offset)
	return e.client.NormalTxByAddress(address, &startBlock, &endBlock, page, offset, false)
}

func (e *ESClient) getAllTxsByAddress(address string) ([]etherscan.NormalTx, error) {
	var offset = 500
	var page = 1
	var all []etherscan.NormalTx

	for {
		txs, err := e.getTxByAddress(address, page, offset)
		if err != nil {
			return all, err
		}

		all = append(all, txs...)

		// no more next page
		if len(txs) < offset {
			break
		}

		page++
	}

	return all, nil
}

// Validate the fromAddr and toAddr. May take a while if the address has too many transactions.
func (e *ESClient) Connected(fromAddr string, toAddr string) (bool, error) {
	var offset = 500
	var page = 1

	for {
		txs, err := e.getTxByAddress(fromAddr, page, offset)
		if err != nil {
			return false, err
		}

		for i := range txs {
			if strings.EqualFold(txs[i].To,toAddr) {
				return true, nil
			}
		}

		// no more next page
		if len(txs) < offset {
			break
		}

		page++
	}

	return false, nil
}
