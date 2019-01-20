package api

import (
	"encoding/base64"
	"fmt"
	"github.com/perlin-network/safu-go/database"
	"github.com/pkg/errors"
	"log"
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

	go func() {
		log.Println("starts crawling")
		list, err := s.esClient.Crawl(req.ScammerAddress)
		if err != nil {
			log.Println("crawl error:", err)
		}

		if err := s.store.InsertGraph(list...); err != nil {
			log.Println("insert error:", err)
		}

		log.Println("finish crawling")
	}()

	return http.StatusOK, res, nil
}

func (s *service) queryAddress(ctx *requestContext) (int, interface{}, error) {
	var req QueryAddressRequest

	if err := ctx.readJSON(&req); err != nil {
		return http.StatusBadRequest, nil, err
	}

	accountRepScores, err := s.getAccountRepScores(req.TargetAddress)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	if accountRepScores > 30 {
		accountRepScores = 30
	}

	scamReportScores, err := s.getScamReportScores(req.TargetAddress)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
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

func (s *service) getGraph(ctx *requestContext) (int, interface{}, error) {
	graph, err := s.store.GetAllVertices()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	type vertex struct {
		Address string `json:"address"`
		//Parents  []string `json:"parents"`
		Children []string `json:"children"`
	}

	var list []*vertex

	for _, g := range graph {
		v := vertex{
			Address: g.Address,
		}

		for c := range g.Children {
			v.Children = append(v.Children, c)
		}

		list = append(list, &v)
	}

	return http.StatusOK, list, nil
}

func (s *service) allScamReports(ctx *requestContext) (int, interface{}, error) {
	resp := AllScamReportResponse{
		Reports: []*ScamReport{},
	}
	s.store.ForEachReport(func(report *database.Report) error {
		sr := &ScamReport{
			ID:             report.ID,
			Timestamp:      report.Timestamp,
			AccountID:      report.AccountID,
			ScammerAddress: report.ScammerAddress,
			VictimAddress:  report.VictimAddress,
			Title:          report.Title,
			Content:        report.Content,
			Proof:          report.Proof,
			Taint:          report.Taint,
		}
		resp.Reports = append(resp.Reports, sr)
		return nil
	})
	return http.StatusOK, resp, nil
}

//////////////////////////////////////////////

func (s *service) getAccountRepScores(targetAddress string) (int, error) {
	reports, err := s.store.GetReportsByScamAddress(targetAddress)
	if err != nil {
		return 0, err
	}
	var accountIDs []string
	for _, report := range reports {
		accountIDs = append(accountIDs, report.AccountID)
	}
	rep, err := s.ledger.GetReps(accountIDs)
	if err != nil {
		return 0, err
	}
	return rep * 10, nil
}

func (s *service) getScamReportScores(targetAddress string) (int, error) {
	report, err := s.store.GetReportByScamAddress(targetAddress)
	if err != nil {
		return 0, err
	}
	taint := int(report.Taint * 7 / 10)
	if taint > 70 {
		taint = 70
	}
	return taint, nil
}
