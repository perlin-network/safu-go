package etherscan

import (
	"github.com/nanmu42/etherscan-api"
	"github.com/perlin-network/safu-go/model"
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

func (e *ESClient) Crawl(address string) ([]*model.Vertex, error) {
	m, err := e.crawl(address)
	if err != nil {
		return nil, err
	}

	var list []*model.Vertex
	for _, v := range m {
		list = append(list, v)
	}

	return list, nil
}

func (e *ESClient) crawl(address string) (map[string]*model.Vertex, error) {
	address = strings.ToLower(address)

	var graph = make(map[string]*model.Vertex)

	txs, err := e.getTxByAddress(address, 1, 100)
	if err != nil {
		return nil, err
	}

	parent := model.NewVertex(address)
	graph[address] = parent

	visited := make(map[string]bool)
	visited[address] = true

	var q []etherscan.NormalTx
	q = append(q, txs...)

	for len(q) != 0 {
		u := q[0]
		q = q[1:len(q):len(q)]

		toTxs, err := e.getTxByAddress(u.To, 1, 100)
		//log.Printf("q: %d,get: %s, result:= %d", len(q), u.To, len(toTxs))

		if err != nil {
			continue
		}

		for _, toTx := range toTxs {
			// Convert the to and from address into lower case
			toTx.To = strings.ToLower(toTx.To)
			toTx.From = strings.ToLower(toTx.From)

			if toTx.From != u.From {
				continue
			}

			if _, ok := visited[toTx.To]; ok {
				continue
			}

			//log.Printf("u from: %s, to: %s", u.From, u.To)
			//log.Printf("toTx from: %s, to: %s", toTx.From, toTx.To)

			v := model.NewVertex(toTx.To)
			v.Parents[parent.Address] = struct{}{}
			graph[toTx.To] = v

			parent.Children[toTx.To] = struct{}{}

			visited[toTx.To] = true
			q = append(q, toTx)
		}
	}

	return graph, nil
}

// WIP
func (e *ESClient) getTxByAddress(address string, page int, offset int) ([]etherscan.NormalTx, error) {
	startBlock := 0
	endBlock := 99999999

	//log.Printf("getTxByAddress address: %s, page: %d, offset: %d", address, page, offset)
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
