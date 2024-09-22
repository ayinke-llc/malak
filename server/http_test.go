package server

import (
	"testing"
	"time"

	"github.com/adelowo/gulter"
	"github.com/ayinke-llc/malak/internal/pkg/jwttoken"
	"github.com/ayinke-llc/malak/internal/pkg/socialauth"
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

// This does nothing really
// It is only here to verify that the routing works correctly
// e.g middlewares are correctly set, paths are correctly named.
// Will help catch issues with paths like `/updates/{references` that is
// missing an ending brace or wrongly placed middlewares and others
func TestServer_New(t *testing.T) {

	t.Run("without swagger", func(t *testing.T) {
		cfg := getConfig()

		srv, closeFn := New(getLogger(t), cfg, &bun.DB{},
			jwttoken.New(cfg), socialauth.NewGoogle(cfg), &httplimit.Middleware{},
			&gulter.Gulter{})

		closeFn()

		go func() {
			srv.ListenAndServe()
		}()

		time.Sleep(time.Second * 2)

		require.NoError(t, srv.Close())
	})

	t.Run("with swagger enabled", func(t *testing.T) {

		cfg := getConfig()

		cfg.HTTP.Swagger.UIEnabled = true
		cfg.HTTP.Swagger.Port = 9999

		srv, closeFn := New(getLogger(t), cfg, &bun.DB{},
			jwttoken.New(cfg), socialauth.NewGoogle(cfg), &httplimit.Middleware{},
			&gulter.Gulter{})

		closeFn()

		go func() {
			srv.ListenAndServe()
		}()

		time.Sleep(time.Second * 2)

		require.NoError(t, srv.Close())
	})
}
