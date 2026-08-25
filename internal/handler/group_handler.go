package handler

import (
	"net/http"

	"bloomfilter/internal/model"
	"bloomfilter/pkg/httpx"
)

func (s *Server) registerGroupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/groups", s.createGroup)
	mux.HandleFunc("GET /api/groups", s.listGroup)
	mux.HandleFunc("GET /api/groups/{id}", s.getGroup)
	mux.HandleFunc("PUT /api/groups/{id}", s.updateGroup)
	mux.HandleFunc("DELETE /api/groups/{id}", s.deleteGroup)
}

type createGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *Server) createGroup(w http.ResponseWriter, r *http.Request) {
	var req createGroupRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.CreateGroup(model.Group{Name: req.Name, Description: req.Description})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, x)
}

func (s *Server) listGroup(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.GroupFilter{
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListGroups(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getGroup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	x, err := s.svc.GetGroup(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

type updateGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *Server) updateGroup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateGroupRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.UpdateGroup(id, model.Group{Name: req.Name, Description: req.Description})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

func (s *Server) deleteGroup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteGroup(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
