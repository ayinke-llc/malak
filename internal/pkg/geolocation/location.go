package geolocation

import (
	"context"
	"net/netip"
)

type GeolocationService interface {
	FindByIP(context.Context, netip.Addr) (country, city string, err error)
}
