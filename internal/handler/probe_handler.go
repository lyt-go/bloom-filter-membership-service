package handler

import (
	"net/http"

	"bloomfilter/internal/model"
	"bloomfilter/pkg/httpx"
)

func (s *Server) registerProbeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/probes", s.createProbe)
	mux.HandleFunc("GET /api/probes", s.listProbe)
	mux.HandleFunc("GET /api/probes/{id}", s.getProbe)
	mux.HandleFunc("DELETE /api/probes/{id}", s.deleteProbe)
}

type createProbeRequest struct {
	FilterID      string `json:"filter_id"`
	Value         string `json:"value"`
	Result        string `json:"result"`
	FalsePositive bool   `json:"false_positive"`
}

func (s *Server) createProbe(w http.ResponseWriter, r *http.Request) {
	var req createProbeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.CreateProbe(model.Probe{FilterID: req.FilterID, Value: req.Value, Result: req.Result, FalsePositive: req.FalsePositive})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, x)
}

func (s *Server) listProbe(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ProbeFilter{
		FilterID: r.URL.Query().Get("filter_id"),
		Result:   r.URL.Query().Get("result"),
		Keyword:  r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListProbes(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getProbe(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	x, err := s.svc.GetProbe(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

func (s *Server) deleteProbe(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteProbe(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
