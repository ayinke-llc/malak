package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ayinke-llc/malak"
	malak_mocks "github.com/ayinke-llc/malak/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestWorkspaceHandler_Overview(t *testing.T) {
	for _, v := range generateWorkspaceOverviewTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			updateRepo := malak_mocks.NewMockUpdateRepository(controller)
			deckRepo := malak_mocks.NewMockDeckRepository(controller)
			contactRepo := malak_mocks.NewMockContactRepository(controller)
			shareRepo := malak_mocks.NewMockContactShareRepository(controller)
			fundingRepo := malak_mocks.NewMockFundraisingPipelineRepository(controller)

			v.mockFn(updateRepo, deckRepo, contactRepo, shareRepo, fundingRepo)

			a := &workspaceHandler{
				cfg:         getConfig(),
				updateRepo:  updateRepo,
				deckRepo:    deckRepo,
				contactRepo: contactRepo,
				shareRepo:   shareRepo,
				fundingRepo: fundingRepo,
			}

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{
				ID: workspaceID,
			}))

			WrapMalakHTTPHandler(getLogger(t),
				a.overview, getConfig(), "workspaces.overview").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateWorkspaceOverviewTestTable() []struct {
	name               string
	mockFn             func(updateRepo *malak_mocks.MockUpdateRepository, deckRepo *malak_mocks.MockDeckRepository, contactRepo *malak_mocks.MockContactRepository, shareRepo *malak_mocks.MockContactShareRepository, fundingRepo *malak_mocks.MockFundraisingPipelineRepository)
	expectedStatusCode int
} {
	return []struct {
		name               string
		mockFn             func(updateRepo *malak_mocks.MockUpdateRepository, deckRepo *malak_mocks.MockDeckRepository, contactRepo *malak_mocks.MockContactRepository, shareRepo *malak_mocks.MockContactShareRepository, fundingRepo *malak_mocks.MockFundraisingPipelineRepository)
		expectedStatusCode int
	}{
		{
			name: "update repo fails",
			mockFn: func(updateRepo *malak_mocks.MockUpdateRepository, deckRepo *malak_mocks.MockDeckRepository, contactRepo *malak_mocks.MockContactRepository, shareRepo *malak_mocks.MockContactShareRepository, fundingRepo *malak_mocks.MockFundraisingPipelineRepository) {
				updateRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("could not fetch updates overview"))

				deckRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.DeckOverview{}, nil)

				contactRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.ContactOverview{}, nil)

				shareRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.ShareOverview{}, nil)

				fundingRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.FundingPipelineOverview{}, nil)
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "deck repo fails",
			mockFn: func(updateRepo *malak_mocks.MockUpdateRepository, deckRepo *malak_mocks.MockDeckRepository, contactRepo *malak_mocks.MockContactRepository, shareRepo *malak_mocks.MockContactShareRepository, fundingRepo *malak_mocks.MockFundraisingPipelineRepository) {
				updateRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.UpdateOverview{}, nil)

				deckRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("could not fetch decks overview"))

				contactRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.ContactOverview{}, nil)

				shareRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.ShareOverview{}, nil)

				fundingRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.FundingPipelineOverview{}, nil)
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "contact repo fails",
			mockFn: func(updateRepo *malak_mocks.MockUpdateRepository, deckRepo *malak_mocks.MockDeckRepository, contactRepo *malak_mocks.MockContactRepository, shareRepo *malak_mocks.MockContactShareRepository, fundingRepo *malak_mocks.MockFundraisingPipelineRepository) {
				updateRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.UpdateOverview{}, nil)

				deckRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.DeckOverview{}, nil)

				contactRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("could not fetch contacts overview"))

				shareRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.ShareOverview{}, nil)

				fundingRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.FundingPipelineOverview{}, nil)
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "share repo fails",
			mockFn: func(updateRepo *malak_mocks.MockUpdateRepository, deckRepo *malak_mocks.MockDeckRepository, contactRepo *malak_mocks.MockContactRepository, shareRepo *malak_mocks.MockContactShareRepository, fundingRepo *malak_mocks.MockFundraisingPipelineRepository) {
				updateRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.UpdateOverview{}, nil)

				deckRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.DeckOverview{}, nil)

				contactRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.ContactOverview{}, nil)

				shareRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("could not fetch shares overview"))

				fundingRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.FundingPipelineOverview{}, nil)
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "funding repo fails",
			mockFn: func(updateRepo *malak_mocks.MockUpdateRepository, deckRepo *malak_mocks.MockDeckRepository, contactRepo *malak_mocks.MockContactRepository, shareRepo *malak_mocks.MockContactShareRepository, fundingRepo *malak_mocks.MockFundraisingPipelineRepository) {
				updateRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.UpdateOverview{}, nil)

				deckRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.DeckOverview{}, nil)

				contactRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.ContactOverview{}, nil)

				shareRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.ShareOverview{}, nil)

				fundingRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("could not fetch funding overview"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "all repos succeed",
			mockFn: func(updateRepo *malak_mocks.MockUpdateRepository, deckRepo *malak_mocks.MockDeckRepository, contactRepo *malak_mocks.MockContactRepository, shareRepo *malak_mocks.MockContactShareRepository, fundingRepo *malak_mocks.MockFundraisingPipelineRepository) {
				updateRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.UpdateOverview{}, nil)

				deckRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.DeckOverview{}, nil)

				contactRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.ContactOverview{}, nil)

				shareRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.ShareOverview{}, nil)

				fundingRepo.EXPECT().
					Overview(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.FundingPipelineOverview{}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}
