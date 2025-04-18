package postgres

import (
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/stretchr/testify/require"
)

func TestPlan_List(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	plan := NewPlanRepository(client)

	plans, err := plan.List(t.Context())
	require.NoError(t, err)

	require.Len(t, plans, 2)
}

func TestPlan_Get(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	plan := NewPlanRepository(client)

	_, err := plan.Get(t.Context(), &malak.FetchPlanOptions{
		Reference: "prod_QmtErtydaJZymT",
	})
	require.NoError(t, err)

	_, err = plan.Get(t.Context(), &malak.FetchPlanOptions{
		Reference: "prod_QmtErtyda",
	})
	require.Error(t, err)
	require.Equal(t, malak.ErrPlanNotFound, err)
}

func TestPlan_SetDefault(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	planRepo := NewPlanRepository(client)

	plan, err := planRepo.Get(t.Context(), &malak.FetchPlanOptions{
		Reference: "prod_QmtErtydaJZymT",
	})
	require.NoError(t, err)
	require.NotNil(t, plan)
	require.False(t, plan.IsDefault)

	secondPlan, err := planRepo.Get(t.Context(), &malak.FetchPlanOptions{
		Reference: "prod_QmtFLR9JvXLryD",
	})
	require.NoError(t, err)
	require.NotNil(t, secondPlan)
	require.False(t, secondPlan.IsDefault)

	require.NoError(t, planRepo.SetDefault(t.Context(), plan))

	plan1FromDB, err := planRepo.Get(t.Context(), &malak.FetchPlanOptions{
		Reference: plan.Reference,
	})
	require.NoError(t, err)
	require.NotNil(t, plan)
	require.True(t, plan1FromDB.IsDefault)
}
