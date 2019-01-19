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
		VictimAddress: req.VictimAddress,
		Title: req.Title,
		Content: req.Content,
		Proof: req.Proof,
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

	return http.StatusOK, res, nil
}

func (s *service) queryAddress(ctx *requestContext) (int, interface{}, error) {
	var req QueryAddressRequest

	if err := ctx.readJSON(&req); err != nil {
		return http.StatusBadRequest, nil, err
	}

	// TODO: acutally replace this
	var res = QueryAddressResponse{
		TargetAddress: req.TargetAddress,
		TaintScore:    int32(req.Timestamp % 100),
	}

	return http.StatusOK, res, nil
}
