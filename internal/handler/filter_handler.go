package handler

import (
	"net/http"

	"bloomfilter/internal/model"
	"bloomfilter/pkg/httpx"
)

func (s *Server) registerFilterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/filters", s.createFilter)
	mux.HandleFunc("GET /api/filters", s.listFilter)
	mux.HandleFunc("GET /api/filters/{id}", s.getFilter)
	mux.HandleFunc("PUT /api/filters/{id}", s.updateFilter)
	mux.HandleFunc("DELETE /api/filters/{id}", s.deleteFilter)
	mux.HandleFunc("PATCH /api/filters/{id}/status", s.transitionFilter)
}

type createFilterRequest struct {
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	BitSize           int     `json:"bit_size"`
	HashCount         int     `json:"hash_count"`
	FalsePositiveRate float64 `json:"false_positive_rate"`
	StrategyID        string  `json:"strategy_id"`
	GroupID           string  `json:"group_id"`
}

func (s *Server) createFilter(w http.ResponseWriter, r *http.Request) {
	var req createFilterRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.CreateFilter(model.Filter{
		Name:              req.Name,
		Description:       req.Description,
		BitSize:           req.BitSize,
		HashCount:         req.HashCount,
		FalsePositiveRate: req.FalsePositiveRate,
		StrategyID:        req.StrategyID,
		GroupID:           req.GroupID,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, x)
}

func (s *Server) listFilter(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.FilterFilter{
		Status:     r.URL.Query().Get("status"),
		StrategyID: r.URL.Query().Get("strategy_id"),
		GroupID:    r.URL.Query().Get("group_id"),
		Keyword:    r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListFilters(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getFilter(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	x, err := s.svc.GetFilter(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

type updateFilterRequest struct {
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	BitSize           int     `json:"bit_size"`
	HashCount         int     `json:"hash_count"`
	FalsePositiveRate float64 `json:"false_positive_rate"`
	StrategyID        string  `json:"strategy_id"`
	GroupID           string  `json:"group_id"`
}

func (s *Server) updateFilter(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateFilterRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.UpdateFilter(id, model.Filter{
		Name:              req.Name,
		Description:       req.Description,
		BitSize:           req.BitSize,
		HashCount:         req.HashCount,
		FalsePositiveRate: req.FalsePositiveRate,
		StrategyID:        req.StrategyID,
		GroupID:           req.GroupID,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

func (s *Server) deleteFilter(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteFilter(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type transitionFilterRequest struct {
	Status string `json:"status"`
}

func (s *Server) transitionFilter(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req transitionFilterRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.TransitionFilter(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}
