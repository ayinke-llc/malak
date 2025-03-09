package maxmind

import (
	"context"
	"net/netip"

	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/geolocation"
	"github.com/oschwald/maxminddb-golang/v2"
	"golang.org/x/sync/errgroup"
)

type maxMindImpl struct {
	cityDB    *maxminddb.Reader
	countryDB *maxminddb.Reader
}

func New(cfg config.Config) (geolocation.GeolocationService, error) {
	cityDB, err := maxminddb.Open(cfg.Analytics.MaxMindCityDB)
	if err != nil {
		return nil, err
	}

	countryDB, err := maxminddb.Open(cfg.Analytics.MaxMindCountryDB)
	if err != nil {
		return nil, err
	}

	return &maxMindImpl{
		cityDB:    cityDB,
		countryDB: countryDB,
	}, nil
}

func (m *maxMindImpl) FindByIP(ctx context.Context, addr netip.Addr) (string, string, error) {
	var g errgroup.Group
	var country, city string

	g.Go(func() error {
		var record map[string]any
		if err := m.countryDB.Lookup(addr).Decode(&record); err != nil {
			return err
		}

		if countryData, ok := record["country"].(map[string]any); ok {
			if names, ok := countryData["names"].(map[string]any); ok {
				if name, ok := names["en"].(string); ok {
					country = name
				}
			}
		}
		return nil
	})

	g.Go(func() error {
		var record map[string]any
		if err := m.cityDB.Lookup(addr).Decode(&record); err != nil {
			return err
		}

		if cityData, ok := record["city"].(map[string]any); ok {
			if names, ok := cityData["names"].(map[string]any); ok {
				if name, ok := names["en"].(string); ok {
					city = name
				}
			}
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return "", "", err
	}

	return country, city, nil
}
