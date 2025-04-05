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

		result, total, err := fundingRepo.List(t.Context(), malak.ListPipelineOptions{
			WorkspaceID: workspace.ID,
			Paginator: malak.Paginator{
				Page:    1,
				PerPage: 10,
			},
		})
		require.NoError(t, err)
		require.GreaterOrEqual(t, total, int64(1))
		found := false
		for _, p := range result {
			if p.ID == pipeline.ID {
				found = true
				break
			}
		}
		require.True(t, found, "newly created pipeline should be found in list results")

		var columns []malak.FundraisingPipelineColumn
		err = client.NewSelect().
			Model(&columns).
			Where("fundraising_pipeline_id = ?", pipeline.ID).
			Scan(t.Context())
		require.NoError(t, err)
		require.Len(t, columns, 0)
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
				Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipelineColumn),
			},
			{
				Title:       "Custom Column 2",
				ColumnType:  malak.FundraisePipelineColumnTypeNormal,
				Description: "Custom column description 2",
				Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipelineColumn),
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
		require.Len(t, columns, 2)

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

func TestFundraising_List(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	fundingRepo := NewFundingRepo(client)
	workspaceRepo := NewWorkspaceRepository(client)

	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	pipelines := []malak.FundraisingPipeline{
		{
			Reference:         malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipeline),
			WorkspaceID:       workspace.ID,
			Title:             "Pipeline 1",
			Stage:             malak.FundraisePipelineStageSeed,
			Description:       "Description 1",
			TargetAmount:      1000000,
			StartDate:         time.Now().UTC(),
			ExpectedCloseDate: time.Now().UTC().Add(90 * 24 * time.Hour),
			IsClosed:          false,
		},
		{
			Reference:         malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipeline),
			WorkspaceID:       workspace.ID,
			Title:             "Pipeline 2",
			Stage:             malak.FundraisePipelineStageSeriesA,
			Description:       "Description 2",
			TargetAmount:      5000000,
			StartDate:         time.Now().UTC(),
			ExpectedCloseDate: time.Now().UTC().Add(90 * 24 * time.Hour),
			IsClosed:          true,
		},
		{
			Reference:         malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipeline),
			WorkspaceID:       workspace.ID,
			Title:             "Pipeline 3",
			Stage:             malak.FundraisePipelineStageBridgeRound,
			Description:       "Description 3",
			TargetAmount:      2000000,
			StartDate:         time.Now().UTC(),
			ExpectedCloseDate: time.Now().UTC().Add(90 * 24 * time.Hour),
			IsClosed:          false,
		},
	}

	for i := range pipelines {
		err = fundingRepo.Create(t.Context(), &pipelines[i])
		require.NoError(t, err)
	}

	t.Run("list all pipelines with pagination", func(t *testing.T) {
		opts := malak.ListPipelineOptions{
			WorkspaceID: workspace.ID,
			Paginator: malak.Paginator{
				Page:    1,
				PerPage: 2,
			},
		}

		result, total, err := fundingRepo.List(t.Context(), opts)
		require.NoError(t, err)
		require.Equal(t, int64(3), total)
		require.Len(t, result, 2)

		opts.Paginator.Page = 2
		result, total, err = fundingRepo.List(t.Context(), opts)
		require.NoError(t, err)
		require.Equal(t, int64(3), total)
		require.Len(t, result, 1)
	})

	t.Run("list active pipelines only", func(t *testing.T) {
		opts := malak.ListPipelineOptions{
			WorkspaceID: workspace.ID,
			ActiveOnly:  true,
			Paginator: malak.Paginator{
				Page:    1,
				PerPage: 10,
			},
		}

		result, total, err := fundingRepo.List(t.Context(), opts)
		require.NoError(t, err)
		require.Equal(t, int64(2), total)
		require.Len(t, result, 2)

		for _, pipeline := range result {
			require.False(t, pipeline.IsClosed)
		}
	})

	t.Run("list pipelines for different workspace", func(t *testing.T) {
		opts := malak.ListPipelineOptions{
			WorkspaceID: uuid.New(), // Different workspace ID
			Paginator: malak.Paginator{
				Page:    1,
				PerPage: 10,
			},
		}

		result, total, err := fundingRepo.List(t.Context(), opts)
		require.NoError(t, err)
		require.Equal(t, int64(0), total)
		require.Len(t, result, 0)
	})

	t.Run("verify sort order by created_at DESC", func(t *testing.T) {
		opts := malak.ListPipelineOptions{
			WorkspaceID: workspace.ID,
			Paginator: malak.Paginator{
				Page:    1,
				PerPage: 10,
			},
		}

		result, _, err := fundingRepo.List(t.Context(), opts)
		require.NoError(t, err)
		require.Len(t, result, 3)

		for i := 1; i < len(result); i++ {
			require.True(t, result[i-1].CreatedAt.After(result[i].CreatedAt) ||
				result[i-1].CreatedAt.Equal(result[i].CreatedAt))
		}
	})
}
