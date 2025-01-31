package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/stretchr/testify/require"
)

func TestRequireWorkspaceValidSubscription(t *testing.T) {
	t.Run("sub not active", func(t *testing.T) {

		rr := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Add("Content-Type", "application/json")
		req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

		requireWorkspaceValidSubscription(getConfig())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.NoError(t, json.NewEncoder(w).Encode("{}"))
		})).ServeHTTP(rr, req)

		require.Equal(t, http.StatusPaymentRequired, rr.Code)
		verifyMatch(t, rr)
	})

	t.Run("sub active", func(t *testing.T) {

		rr := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Add("Content-Type", "application/json")
		req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{
			IsSubscriptionActive: true,
		}))

		requireWorkspaceValidSubscription(getConfig())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.NoError(t, json.NewEncoder(w).Encode("{}"))
		})).ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		verifyMatch(t, rr)
	})
}
