package database

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/perlin-network/safu-go/model"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"strings"
)

type TieDotStore struct {
	db *leveldb.DB
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

	b, err := json.Marshal(report)
	if err != nil {
		return "", err
	}

	if err := t.db.Put([]byte(key), b, nil); err != nil {
		return "", err
	}

	index := "report_index_" + report.ScammerAddress
	if err := t.db.Put([]byte(index), []byte(key), nil); err != nil {
		return "", err
	}

	return id.String(), nil
}

func (t *TieDotStore) GetReport(scammerAddress string) (*Report, error) {
	index := "report_index_" + scammerAddress

	key, err := t.db.Get([]byte(index), nil)
	if err != nil {
		return nil, err
	}

	b, err := t.db.Get(key, nil)
	if err != nil {
		return nil, err
	}

	var report Report

	err = json.Unmarshal(b, &report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (t *TieDotStore) InsertGraph(graph []*model.Vertex) error {
	for _, v := range graph {

		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		key := "addr_" + v.Address
		err = t.db.Put([]byte(key), b, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TieDotStore) BFS(address string) error {
	address = strings.ToLower(address)

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

		log.Println(v.String())

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

func (t *TieDotStore) Close() {
	_ = t.db.Close()
}
