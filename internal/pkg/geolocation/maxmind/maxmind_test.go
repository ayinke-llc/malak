package maxmind

import (
	"context"
	"net/netip"
	"testing"

	"github.com/ayinke-llc/malak/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaxMind_FindByIP(t *testing.T) {
	tests := []struct {
		name        string
		ip          string
		wantCountry string
		wantCity    string
		wantErr     bool
	}{
		{
			name:        "Valid IP - United States/New York",
			ip:          "38.132.98.195",
			wantCountry: "United States",
			wantCity:    "New York",
			wantErr:     false,
		},
		{
			name:        "Invalid IP",
			ip:          "invalid",
			wantCountry: "",
			wantCity:    "",
			wantErr:     true,
		},
	}

	cfg := config.Config{
		Analytics: struct {
			MaxMindCountryDB string "json:\"max_mind_country_db,omitempty\" yaml:\"max_mind_country_db\" mapstructure:\"max_mind_country_db\""
			MaxMindCityDB    string "json:\"max_mind_city_db,omitempty\" yaml:\"max_mind_city_db\" mapstructure:\"max_mind_city_db\""
		}{
			MaxMindCityDB:    "testdata/city.mmdb",
			MaxMindCountryDB: "testdata/country.mmdb",
		},
	}

	service, err := New(cfg)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := netip.ParseAddr(tt.ip)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			country, city, err := service.FindByIP(context.Background(), addr)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCountry, country)
			assert.Equal(t, tt.wantCity, city)
		})
	}
}
