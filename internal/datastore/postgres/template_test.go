package postgres

import (
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/stretchr/testify/require"
)

func TestTemplateRepository_System(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	templateRepo := NewTemplateRepository(client)

	t.Run("returns default templates ordered by updated_at DESC", func(t *testing.T) {
		templates, err := templateRepo.System(t.Context(), malak.SystemTemplateFilterAll)
		require.NoError(t, err)
		require.Len(t, templates, 2)
		require.Equal(t, "500 co test item", templates[0].Title)
		require.Equal(t, "Oops template", templates[1].Title)
	})

	t.Run("returns templates ordered by most used", func(t *testing.T) {
		// Test most used filter on default templates
		result, err := templateRepo.System(t.Context(), malak.SystemTemplateFilterMostUsed)
		require.NoError(t, err)
		require.Len(t, result, 2)
		// From fixtures: "Oops template" has 1 use, "500 co test item" has 0 uses
		require.Equal(t, "Oops template", result[0].Title)
		require.Equal(t, "500 co test item", result[1].Title)
	})

	t.Run("returns templates ordered by recently created", func(t *testing.T) {
		// Test recently created filter on default templates
		result, err := templateRepo.System(t.Context(), malak.SystemTemplateFilterRecentlyCreated)
		require.NoError(t, err)
		require.Len(t, result, 2)
		// Both templates have same created_at in fixtures
		require.Equal(t, "500 co test item", result[0].Title)
		require.Equal(t, "Oops template", result[1].Title)
	})

	t.Run("handles invalid filter gracefully", func(t *testing.T) {
		// Test with invalid filter
		result, err := templateRepo.System(t.Context(), "invalid_filter")
		require.NoError(t, err)
		// Should default to all templates with default ordering (updated_at DESC)
		require.Len(t, result, 2)
		require.Equal(t, "500 co test item", result[0].Title)
		require.Equal(t, "Oops template", result[1].Title)
	})
}
