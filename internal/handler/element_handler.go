package handler

import (
	"net/http"

	"bloomfilter/internal/model"
	"bloomfilter/pkg/httpx"
)

func (s *Server) registerElementRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/elements", s.createElement)
	mux.HandleFunc("GET /api/elements", s.listElement)
	mux.HandleFunc("GET /api/elements/{id}", s.getElement)
	mux.HandleFunc("DELETE /api/elements/{id}", s.deleteElement)
}

type createElementRequest struct {
	FilterID string `json:"filter_id"`
	Value    string `json:"value"`
	Category string `json:"category"`
}

func (s *Server) createElement(w http.ResponseWriter, r *http.Request) {
	var req createElementRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.CreateElement(model.Element{FilterID: req.FilterID, Value: req.Value, Category: req.Category})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, x)
}

func (s *Server) listElement(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ElementFilter{
		FilterID: r.URL.Query().Get("filter_id"),
		Category: r.URL.Query().Get("category"),
		Keyword:  r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListElements(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getElement(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	x, err := s.svc.GetElement(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

func (s *Server) deleteElement(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteElement(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
