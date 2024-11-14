package server

import (
	"net/http"

	"github.com/ayinke-llc/malak"
	"github.com/go-chi/render"
)

type GenericRequest struct{}

func (g GenericRequest) Bind(_ *http.Request) error { return nil }

type meta struct {
	Paging pagingInfo `json:"paging,omitempty" validate:"required"`
}

type pagingInfo struct {
	Total   int64 `json:"total,omitempty" validate:"required"`
	PerPage int64 `json:"per_page,omitempty" validate:"required"`
	Page    int64 `json:"page,omitempty" validate:"required"`
}

type APIStatus struct {
	statusCode int
	Message    string `json:"message,omitempty" validate:"required"`
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
	User             malak.User        `json:"user,omitempty" validate:"required"`
	Workspaces       []malak.Workspace `json:"workspaces,omitempty" validate:"required"`
	CurrentWorkspace *malak.Workspace  `json:"current_workspace,omitempty" validate:"optional"`
	Token            string            `json:"token,omitempty" validate:"required"`
	APIStatus
}

type fetchWorkspaceResponse struct {
	Workspace malak.Workspace `json:"workspace,omitempty" validate:"required"`
	APIStatus `validate:"required"`
}

type fetchContactResponse struct {
	Contact malak.Contact `json:"contact,omitempty" validate:"required"`
	APIStatus
}

type createdUpdateResponse struct {
	Update malak.Update `json:"update,omitempty" validate:"required"`
	APIStatus
}

type listUpdateResponse struct {
	Updates []malak.Update `json:"updates,omitempty" validate:"required"`
	Meta    meta           `json:"meta,omitempty" validate:"required"`
	APIStatus
}

type uploadImageResponse struct {
	URL string `json:"url,omitempty" validate:"required"`
	APIStatus
}

type fetchUpdateReponse struct {
	Update malak.Update `json:"update,omitempty" validate:"required"`
	APIStatus
}

type fetchContactListResponse struct {
	APIStatus
	List malak.ContactList `json:"list,omitempty"`
}

type fetchContactListsResponse struct {
	APIStatus
	Lists []struct {
		List     malak.ContactList                     `json:"list,omitempty" validate:"required"`
		Mappings []malak.ContactListMappingWithContact `json:"mappings,omitempty" validate:"required"`
	} `json:"lists,omitempty" validate:"required"`
}
