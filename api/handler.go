package api

import (
	"github.com/pkg/errors"
	"net/http"
)

func (s *service) postScamReport(ctx *requestContext) (int, interface{}, error) {
	// TODO: store a scam report
	return http.StatusBadRequest, nil, errors.New("not implemented")
}

func (s *service) queryAddress(ctx *requestContext) (int, interface{}, error) {
	// TODO: return the taint result
	return http.StatusBadRequest, nil, errors.New("not implemented")
}
