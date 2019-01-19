package api

import (
	"encoding/base64"
	"fmt"
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

	id, err := s.store.AddReport(req.ScammerAddress, req.VictimAddress, req.Title, req.Content, req.Proof)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	reportID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", id)))
	var res = SubmitReportResponse{
		ID: reportID,
	}

	// TODO: spawn process to update the database with scraped values
	// TODO: add to the list of scam reports the list of accounts that reported it

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
	// TODO:
	return 0
}
