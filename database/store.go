package database

import (
	tiedot "github.com/HouzuoGuo/tiedot/db"
	"log"
)

type TieDotStore struct {
	db *tiedot.DB
}

func NewTieDotStore(dir string) *TieDotStore {
	db, err := tiedot.OpenDB(dir)
	if err != nil {
		panic(err)
	}

	if !db.ColExists("reports") {
		err = db.Create("reports")
		if err != nil {
			log.Panicf("create database error: %s", err)
		}
	}

	err = db.Scrub("reports")
	if err != nil {
		log.Panicf("scrub reports error: %s", err)
	}

	return &TieDotStore{
		db: db,
	}
}

func (t *TieDotStore) AddReport(scammerAddr, victimAddr, title, content, proof string) (int, error) {
	reports := t.db.Use("reports")

	docID, err := reports.Insert(map[string]interface{}{
		"scammer_address": scammerAddr,
		"victim_address":  victimAddr,
		"title":           title,
		"content":         content,
		"proof":           proof,
	})

	if err != nil {
		return 0, err
	}

	return docID, nil
}

func (t *TieDotStore)Close() {
	_ = t.db.Close()
}