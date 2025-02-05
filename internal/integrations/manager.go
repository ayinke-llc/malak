package integrations

import (
	"errors"

	"github.com/ayinke-llc/malak"
)

type IntegrationsManager struct {
	byType map[malak.IntegrationProvider]malak.IntegrationProviderClient
}

func NewManager() *IntegrationsManager {
	return &IntegrationsManager{
		byType: make(map[malak.IntegrationProvider]malak.IntegrationProviderClient),
	}
}

func (m *IntegrationsManager) Add(
	typ malak.IntegrationProvider,
	client malak.IntegrationProviderClient) {

	_, ok := m.byType[typ]
	if ok {
		return
	}

	m.byType[typ] = client
}

func (m *IntegrationsManager) Get(typ malak.IntegrationProvider) (malak.IntegrationProviderClient, error) {
	client, ok := m.byType[typ]
	if ok {
		return client, nil
	}

	return nil, errors.New("client does not exists")
}
