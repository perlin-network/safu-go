package api

import (
	"encoding/base64"
	"fmt"
	"github.com/perlin-network/safu-go/database"
	"github.com/pkg/errors"
	"net/http"
)

func (s *service) postScamReport(ctx *requestContext) (int, interface{}, error) {
	var req SubmitReportRequest

	if err := ctx.readJSON(&req); err != nil {
		return http.StatusBadRequest, nil, err
	}

	if err := validate.Struct(req); err != nil {
		return http.StatusBadRequest, nil, errors.Wrap(err, "invalid request")
	}

	report := database.Report{
		ScammerAddress: req.ScammerAddress,
		VictimAddress:  req.VictimAddress,
		Title:          req.Title,
		Content:        req.Content,
		Proof:          req.Proof,
		Timestamp:      req.Timestamp,
		AccountID:      req.AccountID,
	}

	id, err := s.store.AddReport(report)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	reportID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s", id)))
	var res = SubmitReportResponse{
		ID: reportID,
	}

	// TODO: spawn process to update the database with scraped values
	// TODO: add to the list of scam reports the list of accounts that reported it

	//go func() {
	//	list, err := s.esClient.Crawl(req.ScammerAddress)
	//	if err != nil {
	//		log.Println("crawl error:", err)
	//	}
	//
	//	if err := s.store.InsertGraph(list...); err != nil {
	//		log.Println("insert error:", err)
	//	}
	//	log.Println("insert finished")
	//}()

	return http.StatusOK, res, nil
}

func (s *service) queryAddress(ctx *requestContext) (int, interface{}, error) {
	var req QueryAddressRequest

	if err := ctx.readJSON(&req); err != nil {
		return http.StatusBadRequest, nil, err
	}

	accountRepScores := s.getAccountRepScores(req.TargetAddress)
	if accountRepScores > 30 {
		accountRepScores = 30
	}

	scamReportScores := s.getScamReportScores(req.TargetAddress)
	if scamReportScores > 70 {
		accountRepScores = 70
	}
	taintScore := accountRepScores + scamReportScores

	var res = QueryAddressResponse{
		TargetAddress: req.TargetAddress,
		TaintScore:    int32(taintScore),
	}

	return http.StatusOK, res, nil
}

func (s *service) getAccountRepScores(targetAddress string) int {
	// TODO:
	return 0
}

func (s *service) getScamReportScores(targetAddress string) int {
	s.store.BFS(targetAddress)
	return 0
}
