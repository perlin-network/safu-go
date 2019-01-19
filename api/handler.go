package api

import (
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

	var res = SubmitReportResponse{
		ID: id,
	}

	return http.StatusOK, res, nil
}

func (s *service) queryAddress(ctx *requestContext) (int, interface{}, error) {
	// TODO: return the taint result
	return http.StatusBadRequest, nil, errors.New("not implemented")
}
