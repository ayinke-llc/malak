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

type fetchDetailedContactResponse struct {
	Contact     malak.Contact            `json:"contact,omitempty" validate:"required"`
	SharedItems []malak.ContactShareItem `json:"shared_items,omitempty" validate:"required"`
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

type listIntegrationChartsResponse struct {
	Charts []malak.IntegrationChart `json:"charts,omitempty" validate:"required"`
	APIStatus
}

type listChartDataPointsResponse struct {
	DataPoints []malak.IntegrationDataPoint `json:"data_points,omitempty" validate:"required"`
	APIStatus
}

type listDashboardChartsResponse struct {
	Charts    []malak.DashboardChart         `json:"charts,omitempty" validate:"required"`
	Positions []malak.DashboardChartPosition `json:"positions,omitempty" validate:"required"`
	Dashboard malak.Dashboard                `json:"dashboard,omitempty" validate:"required"`
	Link      malak.DashboardLink            `json:"link,omitempty" validate:"required"`
	APIStatus
}

type listDashboardResponse struct {
	Meta       meta              `json:"meta,omitempty" validate:"required"`
	Dashboards []malak.Dashboard `json:"dashboards,omitempty" validate:"required"`
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
	List malak.ContactList `json:"list,omitempty" validate:"required"`
}

type fetchContactListsResponse struct {
	APIStatus
	Lists []struct {
		List     malak.ContactList                     `json:"list,omitempty" validate:"required"`
		Mappings []malak.ContactListMappingWithContact `json:"mappings,omitempty" validate:"required"`
	} `json:"lists,omitempty" validate:"required"`
}

type fetchUpdateAnalyticsResponse struct {
	APIStatus
	Update     malak.UpdateStat        `json:"update,omitempty" validate:"required"`
	Recipients []malak.UpdateRecipient `json:"recipients,omitempty" validate:"required"`
}

type fetchPublicDeckResponse struct {
	Deck malak.PublicDeck `json:"deck,omitempty" validate:"required"`
	APIStatus
}

type fetchDeckResponse struct {
	Deck malak.Deck `json:"deck,omitempty" validate:"required"`
	APIStatus
}

type fetchDecksResponse struct {
	Decks []malak.Deck `json:"decks,omitempty" validate:"required"`
	APIStatus
}

type listContactsResponse struct {
	Contacts []malak.Contact `json:"contacts,omitempty" validate:"required"`
	Meta     meta            `json:"meta,omitempty" validate:"required"`
	APIStatus
}

type preferenceResponse struct {
	Preferences *malak.Preference `json:"preferences,omitempty" validate:"required"`
	APIStatus
}

type listIntegrationResponse struct {
	Integrations []malak.WorkspaceIntegration `json:"integrations,omitempty" validate:"required"`
	APIStatus
}

type fetchBillingPortalResponse struct {
	Link string `json:"link,omitempty" validate:"required"`
	APIStatus
}

type fetchDashboardResponse struct {
	Dashboard malak.Dashboard `json:"dashboard,omitempty" validate:"required"`
	APIStatus
}

type fetchTemplatesResponse struct {
	Templates struct {
		System    []malak.SystemTemplate `json:"system,omitempty" validate:"required"`
		Workspace []malak.SystemTemplate `json:"workspace,omitempty" validate:"required"`
	} `json:"templates,omitempty" validate:"required"`
	APIStatus
}

type fetchSessionsDeck struct {
	Sessions []*malak.DeckViewerSession `json:"sessions,omitempty" validate:"required"`
	Meta     meta                       `json:"meta,omitempty" validate:"required"`
	APIStatus
}

type regenerateLinkResponse struct {
	Link malak.DashboardLink `json:"link,omitempty" validate:"required"`
	APIStatus
}

type listDashboardLinkResponse struct {
	Meta  meta                  `json:"meta,omitempty" validate:"required"`
	Links []malak.DashboardLink `json:"links,omitempty" validate:"required"`
	APIStatus
}

type createdAPIKeyResponse struct {
	APIStatus
	Value string `json:"value,omitempty" validate:"required"`
}
