package server

import (
	"net/http"

	"github.com/ayinke-llc/malak"
	"github.com/go-chi/render"
)

type GenericRequest struct{}

func (g GenericRequest) Bind(_ *http.Request) error { return nil }

type meta struct {
	Paging pagingInfo `json:"paging"`
}

type pagingInfo struct {
	Total   int64 `json:"total"`
	PerPage int64 `json:"per_page"`
	Page    int64 `json:"page"`
}

type APIStatus struct {
	statusCode int
	// Generic message that tells you the status of the operation
	Message string `json:"message"`
}

func (a APIStatus) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, a.statusCode)
	return nil
}

type APIError struct {
	APIStatus
}

func newAPIStatus(code int, s string) APIStatus {
	return APIStatus{
		statusCode: code,
		Message:    s,
	}
}

type createdUserResponse struct {
	User  *malak.User `json:"user"`
	Token string      `json:"token"`
	APIStatus
}

type fetchWorkspaceResponse struct {
	Workspace *malak.Workspace `json:"workspace,omitempty"`
	APIStatus
}
