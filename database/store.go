package database

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/perlin-network/safu-go/model"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type TieDotStore struct {
	db *leveldb.DB
}

type EthTransaction struct {
	Hash string `json:"hash"`
	From string `json:"from"`
	To   string `json:"to"`
}

func NewTieDotStore(dir string) *TieDotStore {
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		log.Panicf("open db error: %s", err)
	}

	return &TieDotStore{
		db: db,
	}
}

func (t *TieDotStore) AddReport(report Report) (string, error) {
	id, _ := uuid.NewV4()

	key := "report_" + id.String()

	report.ID = id.String()

	b, err := json.Marshal(report)
	if err != nil {
		return "", err
	}

	if err := t.db.Put([]byte(key), b, nil); err != nil {
		return "", err
	}

	return id.String(), nil
}

func (t *TieDotStore) GetReportsByScamAddress(scammerAddress string) ([]*Report, error) {
	var reports []*Report
	iter := t.db.NewIterator(nil, nil)
	iter.Release()

	for iter.Next() {
		value := iter.Value()
		var r = &Report{}
		err := json.Unmarshal(value, r)
		if err != nil {
			return nil, err
		}

		if strings.EqualFold(scammerAddress, r.ScammerAddress) {
			reports = append(reports, r)
		}
	}

	return reports, nil
}

func (t *TieDotStore) InsertGraph(graph ...*model.Vertex) error {
	batch := &leveldb.Batch{}
	for _, v := range graph {

		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		key := "addr_" + v.Address

		batch.Put([]byte(key), b)
	}

	return t.db.Write(batch, nil)
}

func (t *TieDotStore) GetReportByScamAddress(address string) (*Report, error) {
	reports, err := t.GetReportsByScamAddress(address)
	if err != nil {
		return nil, err
	}

	if len(reports) < 1 {
		return nil, errors.Errorf("could not find report with scam address %s", address)
	}
	return reports[0], nil
}

func (t *TieDotStore) TaintBFS(address string, taint int) error {
	address = strings.ToLower(address)

	err := t.updateReportsTaints(address, taint)
	if err != nil {
		return err
	}

	childrenTaint := int(0.3 * float32(taint))

	q := []string{address}
	visited := make(map[string]struct{})
	visited[address] = struct{}{}

	for len(q) != 0 {
		u := q[0]
		q = q[1:len(q):len(q)]

		v, err := t.getVertex(u)
		if err != nil {
			return err
		}

		err = t.updateReportsTaints(v.Address, childrenTaint)
		if err != nil {
			return err
		}

		for p := range v.Children {
			if _, ok := visited[p]; !ok {

				q = append(q, p)
				visited[p] = struct{}{}
			}
		}
	}

	return nil
}

func (t *TieDotStore) getVertex(address string) (*model.Vertex, error) {
	key := "addr_" + address

	b, err := t.db.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}

	var v model.Vertex
	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (t *TieDotStore) GetAllVertices() ([]*model.Vertex, error) {
	var list []*model.Vertex
	prefix := "addr_"
	//prefixLength := len(prefix)

	iter := t.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)

	var err error
	for iter.Next() {
		value := iter.Value()

		//log.Println("check: ", string(value))
		var vertex = &model.Vertex{}
		err = json.Unmarshal(value, vertex)
		if err != nil {
			return nil, err
		}

		list = append(list, vertex)
	}

	iter.Release()

	if err != nil {
		return nil, err
	}

	return list, iter.Error()
}

func (t *TieDotStore) updateReportsTaints(address string, taint int) error {
	reports, err := t.GetReportsByScamAddress(address)
	log.Printf("updateReportsTaints len: %d", len(reports))
	if err != nil {
		return err
	}

	batch := &leveldb.Batch{}
	for _, report := range reports {
		report.Taint = taint

		b, err := json.Marshal(report)
		if err != nil {
			return err
		}

		log.Printf("update taint ID: %s, taint: %d", report.ID, report.Taint)
		batch.Put([]byte(report.ID), b)
	}

	return t.db.Write(batch, nil)
}

func (t *TieDotStore) ForEachReport(callback func(report *Report) error) error {
	prefix := "report_"
	//prefixLength := len(prefix)

	iter := t.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)

	var err error
	for iter.Next() {
		value := iter.Value()

		//log.Println("check: ", string(value))
		var r = &Report{}
		err = json.Unmarshal(value, r)
		if err != nil {
			return err
		}

		if err := callback(r); err != nil {
			return err
		}
	}

	iter.Release()

	if err != nil {
		return err
	}

	return iter.Error()
}

func (t *TieDotStore) BFSEthAccount(accountID string, apiKey string) error {
	accountID = strings.ToLower(accountID)

	type APIResponse struct {
		Status  string            `json:"status"`
		Message string            `json:"message"`
		Result  []*EthTransaction `json:"result"`
	}
	ids := list.New()
	ids.PushFront(accountID)
	searched := make(map[string]struct{})
	client := &http.Client{}

	count := 0

	for ids.Len() > 0 {
		if count == 30 {
			break
		}
		count++
		time.Sleep(time.Millisecond * 250)

		accountID = ids.Front().Value.(string)
		ids.Remove(ids.Front())

		if len(searched) >= 256 {
			return nil
		}
		if _, ok := searched[accountID]; ok {
			continue
		}
		searched[accountID] = struct{}{}

		resp, err := client.Get(fmt.Sprintf(
			"http://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&sort=desc&apikey=%s",
			accountID, apiKey,
		))
		if err != nil {
			log.Printf("unable to retrieve transaction list: %+v\n", err)
			continue
		}

		defer resp.Body.Close()
		_out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("unable to decode output: %+v\n", err)
			continue
		}

		var out APIResponse
		err = json.Unmarshal(_out, &out)
		if err != nil {
			log.Printf("unable to unmarshal output: %+v\n", err)
			continue
		}

		for _, tx := range out.Result {
			tx.From = strings.ToLower(tx.From)
			tx.To = strings.ToLower(tx.To)
			if tx.To != accountID {
				continue
			}
			log.Println(tx)

			ids.PushBack(tx.From)
		}
	}

	return nil
}

func (t *TieDotStore) Close() {
	_ = t.db.Close()
}
