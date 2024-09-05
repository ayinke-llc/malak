package postgres

import (
	"context"
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/stretchr/testify/require"
)

func TestPlan_List(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	plan := NewPlanRepository(client)

	plans, err := plan.List(context.Background())
	require.NoError(t, err)

	require.Len(t, plans, 2)
}

func TestPlan_Get(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	plan := NewPlanRepository(client)

	_, err := plan.Get(context.Background(), &malak.FetchPlanOptions{
		Reference: "prod_QmtErtydaJZymT",
	})
	require.NoError(t, err)
}
