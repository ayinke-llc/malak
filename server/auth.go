package server

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type authHandler struct {
	logger *logrus.Entry
	// googleCfg
}

type authenticateUserRequest struct {
	genericRequest

	Code string `json:"code,omitempty"`
}

// @Summary Sign in with a social login provider
// @Tags auth
// @Accept  json
// @Produce  json
// @Param message body authenticateUserRequest true "auth exchange data"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /auth/login [post]
func (a *authHandler) Login(w http.ResponseWriter, r *http.Request) {

}
