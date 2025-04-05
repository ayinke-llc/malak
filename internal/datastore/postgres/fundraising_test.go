package postgres

import (
	"testing"
	"time"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestFundraising_Create(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	fundingRepo := NewFundingRepo(client)
	workspaceRepo := NewWorkspaceRepository(client)

	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	t.Run("create with default columns", func(t *testing.T) {
		pipeline := &malak.FundraisingPipeline{
			Reference:         malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipeline),
			WorkspaceID:       workspace.ID,
			Title:             "Test Fundraising Round",
			Stage:             malak.FundraisePipelineStageSeed,
			Description:       "Test fundraising round description",
			TargetAmount:      1000000,
			StartDate:         time.Now().UTC(),
			ExpectedCloseDate: time.Now().UTC().Add(90 * 24 * time.Hour),
		}

		err = fundingRepo.Create(t.Context(), pipeline)
		require.NoError(t, err)
		require.NotEmpty(t, pipeline.ID)

		var count int
		count, err = client.NewSelect().
			Model((*malak.FundraisingPipeline)(nil)).
			Where("id = ?", pipeline.ID).
			Count(t.Context())
		require.NoError(t, err)
		require.Equal(t, 1, count)

		var columns []malak.FundraisingPipelineColumn
		err = client.NewSelect().
			Model(&columns).
			Where("fundraising_pipeline_id = ?", pipeline.ID).
			Scan(t.Context())
		require.NoError(t, err)
		require.Len(t, columns, len(malak.DefaultFundraisingColumns))

		columnMap := make(map[string]malak.FundraisingPipelineColumn)
		for _, col := range columns {
			columnMap[col.Title] = col
		}

		for _, defaultCol := range malak.DefaultFundraisingColumns {
			col, exists := columnMap[defaultCol.Title]
			require.True(t, exists)
			require.Equal(t, defaultCol.ColumnType, col.ColumnType)
			require.Equal(t, defaultCol.Description, col.Description)
		}
	})

	t.Run("create with additional columns", func(t *testing.T) {
		pipeline := &malak.FundraisingPipeline{
			Reference:         malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipeline),
			WorkspaceID:       workspace.ID,
			Title:             "Test Fundraising Round with Custom Columns",
			Stage:             malak.FundraisePipelineStageSeed,
			Description:       "Test fundraising round description",
			TargetAmount:      1000000,
			StartDate:         time.Now().UTC(),
			ExpectedCloseDate: time.Now().UTC().Add(90 * 24 * time.Hour),
		}

		additionalColumns := []malak.FundraisingPipelineColumn{
			{
				Title:       "Custom Column 1",
				ColumnType:  malak.FundraisePipelineColumnTypeNormal,
				Description: "Custom column description 1",
			},
			{
				Title:       "Custom Column 2",
				ColumnType:  malak.FundraisePipelineColumnTypeNormal,
				Description: "Custom column description 2",
			},
		}

		err = fundingRepo.Create(t.Context(), pipeline, additionalColumns...)
		require.NoError(t, err)
		require.NotEmpty(t, pipeline.ID)

		var columns []malak.FundraisingPipelineColumn
		err = client.NewSelect().
			Model(&columns).
			Where("fundraising_pipeline_id = ?", pipeline.ID).
			Scan(t.Context())
		require.NoError(t, err)
		require.Len(t, columns, len(malak.DefaultFundraisingColumns)+len(additionalColumns))

		columnMap := make(map[string]malak.FundraisingPipelineColumn)
		for _, col := range columns {
			columnMap[col.Title] = col
		}

		for _, customCol := range additionalColumns {
			col, exists := columnMap[customCol.Title]
			require.True(t, exists)
			require.Equal(t, customCol.ColumnType, col.ColumnType)
			require.Equal(t, customCol.Description, col.Description)
		}
	})
}
