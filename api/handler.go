package api

import (
	"github.com/pkg/errors"
	"net/http"
)

func (s *service) postScamReport(ctx *requestContext) (int, interface{}, error) {
	// TODO: store a scam report
	var req SubmitReportRequest

	if err := ctx.readJSON(&req); err != nil {
		return http.StatusBadRequest, nil, err
	}

	var res SubmitReportResponse

	return http.StatusOK, res, nil
}

func (s *service) queryAddress(ctx *requestContext) (int, interface{}, error) {
	// TODO: return the taint result
	return http.StatusBadRequest, nil, errors.New("not implemented")
}
