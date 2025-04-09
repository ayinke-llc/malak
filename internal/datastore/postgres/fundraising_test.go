package postgres

import (
	"database/sql"
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

func TestFundraising_Get(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	fundingRepo := NewFundingRepo(client)
	workspaceRepo := NewWorkspaceRepository(client)

	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	pipeline := &malak.FundraisingPipeline{
		Reference:         malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipeline),
		WorkspaceID:       workspace.ID,
		Title:             "Test Pipeline",
		Stage:             malak.FundraisePipelineStageSeed,
		Description:       "Test pipeline description",
		TargetAmount:      1000000,
		StartDate:         time.Now().UTC(),
		ExpectedCloseDate: time.Now().UTC().Add(90 * 24 * time.Hour),
	}

	err = fundingRepo.Create(t.Context(), pipeline)
	require.NoError(t, err)

	t.Run("get existing pipeline", func(t *testing.T) {
		result, err := fundingRepo.Get(t.Context(), malak.FetchPipelineOptions{
			WorkspaceID: workspace.ID,
			Reference:   pipeline.Reference,
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, pipeline.ID, result.ID)
		require.Equal(t, pipeline.Title, result.Title)
		require.Equal(t, pipeline.Stage, result.Stage)
	})

	t.Run("get non-existent pipeline", func(t *testing.T) {
		result, err := fundingRepo.Get(t.Context(), malak.FetchPipelineOptions{
			WorkspaceID: workspace.ID,
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipeline),
		})
		require.Error(t, err)
		require.ErrorIs(t, err, malak.ErrPipelineNotFound)
		require.Nil(t, result)
	})

	t.Run("get pipeline from different workspace", func(t *testing.T) {
		result, err := fundingRepo.Get(t.Context(), malak.FetchPipelineOptions{
			WorkspaceID: uuid.New(),
			Reference:   pipeline.Reference,
		})
		require.Error(t, err)
		require.ErrorIs(t, err, malak.ErrPipelineNotFound)
		require.Nil(t, result)
	})
}

func TestFundraising_Board(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	fundingRepo := NewFundingRepo(client)
	workspaceRepo := NewWorkspaceRepository(client)

	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	pipeline := &malak.FundraisingPipeline{
		Reference:         malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipeline),
		WorkspaceID:       workspace.ID,
		Title:             "Test Pipeline",
		Stage:             malak.FundraisePipelineStageSeed,
		Description:       "Test pipeline description",
		TargetAmount:      1000000,
		StartDate:         time.Now().UTC(),
		ExpectedCloseDate: time.Now().UTC().Add(90 * 24 * time.Hour),
	}

	columns := []malak.FundraisingPipelineColumn{
		{
			Title:       "Column 1",
			ColumnType:  malak.FundraisePipelineColumnTypeNormal,
			Description: "First column",
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipelineColumn),
		},
		{
			Title:       "Column 2",
			ColumnType:  malak.FundraisePipelineColumnTypeNormal,
			Description: "Second column",
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipelineColumn),
		},
	}

	err = fundingRepo.Create(t.Context(), pipeline, columns...)
	require.NoError(t, err)

	// Create test contacts first
	testContacts := []*malak.Contact{
		{
			ID:          uuid.New(),
			Email:       malak.Email("test1@example.com"),
			WorkspaceID: workspace.ID,
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeContact),
			FirstName:   "Test",
			LastName:    "Contact1",
		},
		{
			ID:          uuid.New(),
			Email:       malak.Email("test2@example.com"),
			WorkspaceID: workspace.ID,
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeContact),
			FirstName:   "Test",
			LastName:    "Contact2",
		},
	}

	_, err = client.NewInsert().Model(&testContacts).Exec(t.Context())
	require.NoError(t, err)

	contacts := []malak.FundraiseContact{
		{
			ID:                          uuid.New(),
			Reference:                   malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipelineColumnContact),
			ContactID:                   testContacts[0].ID,
			FundraisingPipelineID:       pipeline.ID,
			FundraisingPipelineColumnID: columns[0].ID,
		},
		{
			ID:                          uuid.New(),
			Reference:                   malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipelineColumnContact),
			ContactID:                   testContacts[1].ID,
			FundraisingPipelineID:       pipeline.ID,
			FundraisingPipelineColumnID: columns[1].ID,
		},
	}

	_, err = client.NewInsert().Model(&contacts).Exec(t.Context())
	require.NoError(t, err)

	positions := []malak.FundraiseContactPosition{
		{
			ID:                                 uuid.New(),
			Reference:                          malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipelineColumnContactPosition),
			FundraisingPipelineColumnContactID: contacts[0].ID,
			OrderIndex:                         1,
		},
		{
			ID:                                 uuid.New(),
			Reference:                          malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipelineColumnContactPosition),
			FundraisingPipelineColumnContactID: contacts[1].ID,
			OrderIndex:                         2,
		},
	}

	_, err = client.NewInsert().Model(&positions).Exec(t.Context())
	require.NoError(t, err)

	t.Run("get board with columns, contacts and positions", func(t *testing.T) {
		resultColumns, resultContacts, resultPositions, err := fundingRepo.Board(t.Context(), pipeline)
		require.NoError(t, err)
		require.Len(t, resultColumns, 2)
		require.Len(t, resultContacts, 2)
		require.Len(t, resultPositions, 2)

		for i, col := range resultColumns {
			require.Equal(t, columns[i].Title, col.Title)
			require.Equal(t, columns[i].ColumnType, col.ColumnType)
			require.Equal(t, columns[i].Description, col.Description)
		}

		contactMap := make(map[uuid.UUID]malak.FundraiseContact)
		for _, contact := range resultContacts {
			contactMap[contact.ID] = contact
		}
		for _, expectedContact := range contacts {
			contact, exists := contactMap[expectedContact.ID]
			require.True(t, exists)
			require.Equal(t, expectedContact.ContactID, contact.ContactID)
			require.Equal(t, expectedContact.FundraisingPipelineID, contact.FundraisingPipelineID)
			require.Equal(t, expectedContact.FundraisingPipelineColumnID, contact.FundraisingPipelineColumnID)
		}

		for i := 1; i < len(resultPositions); i++ {
			require.Less(t, resultPositions[i-1].OrderIndex, resultPositions[i].OrderIndex)
		}
	})

	t.Run("get board for pipeline with no columns", func(t *testing.T) {
		emptyPipeline := &malak.FundraisingPipeline{
			Reference:         malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipeline),
			WorkspaceID:       workspace.ID,
			Title:             "Empty Pipeline",
			Stage:             malak.FundraisePipelineStageSeed,
			Description:       "Empty pipeline description",
			TargetAmount:      1000000,
			StartDate:         time.Now().UTC(),
			ExpectedCloseDate: time.Now().UTC().Add(90 * 24 * time.Hour),
		}

		err = fundingRepo.Create(t.Context(), emptyPipeline)
		require.NoError(t, err)

		resultColumns, resultContacts, resultPositions, err := fundingRepo.Board(t.Context(), emptyPipeline)
		require.NoError(t, err)
		require.Empty(t, resultColumns)
		require.Empty(t, resultContacts)
		require.Empty(t, resultPositions)
	})
}

func TestFundraising_CloseBoard(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	fundingRepo := NewFundingRepo(client)
	workspaceRepo := NewWorkspaceRepository(client)

	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	pipeline := &malak.FundraisingPipeline{
		Reference:         malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipeline),
		WorkspaceID:       workspace.ID,
		Title:             "Test Pipeline",
		Stage:             malak.FundraisePipelineStageSeed,
		Description:       "Test pipeline description",
		TargetAmount:      1000000,
		StartDate:         time.Now().UTC(),
		ExpectedCloseDate: time.Now().UTC().Add(90 * 24 * time.Hour),
		IsClosed:          false,
	}

	err = fundingRepo.Create(t.Context(), pipeline)
	require.NoError(t, err)

	t.Run("close board successfully", func(t *testing.T) {
		err := fundingRepo.CloseBoard(t.Context(), pipeline)
		require.NoError(t, err)

		// Verify the pipeline is marked as closed
		result, err := fundingRepo.Get(t.Context(), malak.FetchPipelineOptions{
			WorkspaceID: workspace.ID,
			Reference:   pipeline.Reference,
		})
		require.NoError(t, err)
		require.True(t, result.IsClosed)
	})

	t.Run("close already closed board", func(t *testing.T) {
		err := fundingRepo.CloseBoard(t.Context(), pipeline)
		require.NoError(t, err)

		result, err := fundingRepo.Get(t.Context(), malak.FetchPipelineOptions{
			WorkspaceID: workspace.ID,
			Reference:   pipeline.Reference,
		})
		require.NoError(t, err)
		require.True(t, result.IsClosed)
	})
}

func TestFundraising_AddContactToBoard(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	fundingRepo := NewFundingRepo(client)
	workspaceRepo := NewWorkspaceRepository(client)

	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	pipeline := &malak.FundraisingPipeline{
		Reference:         malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipeline),
		WorkspaceID:       workspace.ID,
		Title:             "Test Pipeline",
		Stage:             malak.FundraisePipelineStageSeed,
		Description:       "Test pipeline description",
		TargetAmount:      1000000,
		StartDate:         time.Now().UTC(),
		ExpectedCloseDate: time.Now().UTC().Add(90 * 24 * time.Hour),
	}

	column := malak.FundraisingPipelineColumn{
		Title:       "Test Column",
		ColumnType:  malak.FundraisePipelineColumnTypeNormal,
		Description: "Test column description",
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipelineColumn),
	}

	err = fundingRepo.Create(t.Context(), pipeline, column)
	require.NoError(t, err)

	contact := &malak.Contact{
		ID:          uuid.New(),
		Email:       malak.Email("test@example.com"),
		WorkspaceID: workspace.ID,
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeContact),
		FirstName:   "Test",
		LastName:    "Contact",
	}

	_, err = client.NewInsert().Model(contact).Exec(t.Context())
	require.NoError(t, err)

	var columns []malak.FundraisingPipelineColumn
	err = client.NewSelect().
		Model(&columns).
		Where("fundraising_pipeline_id = ?", pipeline.ID).
		Scan(t.Context())
	require.NoError(t, err)
	require.Len(t, columns, 1)

	t.Run("successfully add contact to board", func(t *testing.T) {
		opts := &malak.AddContactToBoardOptions{
			Column:             &columns[0],
			Contact:            contact,
			ReferenceGenerator: malak.NewReferenceGenerator(),
		}

		err := fundingRepo.AddContactToBoard(t.Context(), opts)
		require.NoError(t, err)

		var contacts []malak.FundraiseContact
		err = client.NewSelect().
			Model(&contacts).
			Where("fundraising_pipeline_id = ?", pipeline.ID).
			Where("contact_id = ?", contact.ID).
			Scan(t.Context())
		require.NoError(t, err)
		require.Len(t, contacts, 1)
		require.Equal(t, contact.ID, contacts[0].ContactID)
		require.Equal(t, columns[0].ID, contacts[0].FundraisingPipelineColumnID)

		var positions []malak.FundraiseContactPosition
		err = client.NewSelect().
			Model(&positions).
			Where("fundraising_pipeline_column_contact_id = ?", contacts[0].ID).
			Scan(t.Context())
		require.NoError(t, err)
		require.Len(t, positions, 1)
		require.NotZero(t, positions[0].OrderIndex)
	})

	t.Run("add same contact twice should fail", func(t *testing.T) {
		opts := &malak.AddContactToBoardOptions{
			Column:             &columns[0],
			Contact:            contact,
			ReferenceGenerator: malak.NewReferenceGenerator(),
		}

		err := fundingRepo.AddContactToBoard(t.Context(), opts)
		require.Error(t, err) // Should fail due to unique constraint
	})
}

func TestFundraising_DefaultColumn(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	fundingRepo := NewFundingRepo(client)
	workspaceRepo := NewWorkspaceRepository(client)

	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	t.Run("get default column when Backlog column exists", func(t *testing.T) {
		pipeline := &malak.FundraisingPipeline{
			Reference:         malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipeline),
			WorkspaceID:       workspace.ID,
			Title:             "Test Pipeline",
			Stage:             malak.FundraisePipelineStageSeed,
			Description:       "Test pipeline description",
			TargetAmount:      1000000,
			StartDate:         time.Now().UTC(),
			ExpectedCloseDate: time.Now().UTC().Add(90 * 24 * time.Hour),
		}

		columns := []malak.FundraisingPipelineColumn{
			{
				Title:       "First Column",
				ColumnType:  malak.FundraisePipelineColumnTypeNormal,
				Description: "First normal column",
				Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipelineColumn),
			},
			{
				Title:       "Backlog",
				ColumnType:  malak.FundraisePipelineColumnTypeNormal,
				Description: "Default backlog column",
				Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipelineColumn),
			},
			{
				Title:       "Closed Column",
				ColumnType:  malak.FundraisePipelineColumnTypeClosed,
				Description: "Closed column",
				Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipelineColumn),
			},
		}

		err = fundingRepo.Create(t.Context(), pipeline, columns...)
		require.NoError(t, err)

		defaultColumn, err := fundingRepo.DefaultColumn(t.Context(), pipeline)
		require.NoError(t, err)
		require.Equal(t, "Backlog", defaultColumn.Title)
		require.Equal(t, malak.FundraisePipelineColumnTypeNormal, defaultColumn.ColumnType)
	})

	t.Run("get default column when no Backlog column exists", func(t *testing.T) {
		pipeline := &malak.FundraisingPipeline{
			Reference:         malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipeline),
			WorkspaceID:       workspace.ID,
			Title:             "Pipeline Without Backlog",
			Stage:             malak.FundraisePipelineStageSeed,
			Description:       "Pipeline without backlog column",
			TargetAmount:      1000000,
			StartDate:         time.Now().UTC(),
			ExpectedCloseDate: time.Now().UTC().Add(90 * 24 * time.Hour),
		}

		columns := []malak.FundraisingPipelineColumn{
			{
				Title:       "First Column",
				ColumnType:  malak.FundraisePipelineColumnTypeNormal,
				Description: "First normal column",
				Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipelineColumn),
			},
			{
				Title:       "Second Column",
				ColumnType:  malak.FundraisePipelineColumnTypeNormal,
				Description: "Second normal column",
				Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipelineColumn),
			},
		}

		err = fundingRepo.Create(t.Context(), pipeline, columns...)
		require.NoError(t, err)

		_, err := fundingRepo.DefaultColumn(t.Context(), pipeline)
		require.Error(t, err)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("get default column from pipeline with no columns", func(t *testing.T) {
		emptyPipeline := &malak.FundraisingPipeline{
			Reference:         malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipeline),
			WorkspaceID:       workspace.ID,
			Title:             "Empty Pipeline",
			Stage:             malak.FundraisePipelineStageSeed,
			Description:       "Empty pipeline description",
			TargetAmount:      1000000,
			StartDate:         time.Now().UTC(),
			ExpectedCloseDate: time.Now().UTC().Add(90 * 24 * time.Hour),
		}

		err = fundingRepo.Create(t.Context(), emptyPipeline)
		require.NoError(t, err)

		_, err := fundingRepo.DefaultColumn(t.Context(), emptyPipeline)
		require.Error(t, err)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestFundraising_GetContact(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	fundingRepo := NewFundingRepo(client)
	workspaceRepo := NewWorkspaceRepository(client)

	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	pipeline := &malak.FundraisingPipeline{
		Reference:         malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipeline),
		WorkspaceID:       workspace.ID,
		Title:             "Test Pipeline",
		Stage:             malak.FundraisePipelineStageSeed,
		Description:       "Test pipeline description",
		TargetAmount:      1000000,
		StartDate:         time.Now().UTC(),
		ExpectedCloseDate: time.Now().UTC().Add(90 * 24 * time.Hour),
	}

	column := malak.FundraisingPipelineColumn{
		Title:       "Test Column",
		ColumnType:  malak.FundraisePipelineColumnTypeNormal,
		Description: "Test column description",
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipelineColumn),
	}

	err = fundingRepo.Create(t.Context(), pipeline, column)
	require.NoError(t, err)

	contact := &malak.Contact{
		ID:          uuid.New(),
		Email:       malak.Email("test@example.com"),
		WorkspaceID: workspace.ID,
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeContact),
		FirstName:   "Test",
		LastName:    "Contact",
	}

	_, err = client.NewInsert().Model(contact).Exec(t.Context())
	require.NoError(t, err)

	var columns []malak.FundraisingPipelineColumn
	err = client.NewSelect().
		Model(&columns).
		Where("fundraising_pipeline_id = ?", pipeline.ID).
		Scan(t.Context())
	require.NoError(t, err)
	require.Len(t, columns, 1)

	// Add contact to board first
	opts := &malak.AddContactToBoardOptions{
		Column:             &columns[0],
		Contact:            contact,
		ReferenceGenerator: malak.NewReferenceGenerator(),
		Rating:             3,
		CanLeadRound:       true,
		CheckSize:          1000000, // $1M
	}

	err = fundingRepo.AddContactToBoard(t.Context(), opts)
	require.NoError(t, err)

	t.Run("get existing contact", func(t *testing.T) {
		result, err := fundingRepo.GetContact(t.Context(), pipeline.ID, contact.ID)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, contact.ID, result.ContactID)
		require.Equal(t, pipeline.ID, result.FundraisingPipelineID)
		require.Equal(t, columns[0].ID, result.FundraisingPipelineColumnID)
		require.NotNil(t, result.Contact)
		require.Equal(t, contact.FirstName, result.Contact.FirstName)
		require.Equal(t, contact.LastName, result.Contact.LastName)
		require.Equal(t, contact.Email, result.Contact.Email)
	})

	t.Run("get non-existent contact", func(t *testing.T) {
		result, err := fundingRepo.GetContact(t.Context(), pipeline.ID, uuid.New())
		require.Error(t, err)
		require.ErrorIs(t, err, malak.ErrContactNotFoundOnBoard)
		require.Nil(t, result)
	})

	t.Run("get contact from non-existent pipeline", func(t *testing.T) {
		result, err := fundingRepo.GetContact(t.Context(), uuid.New(), contact.ID)
		require.Error(t, err)
		require.ErrorIs(t, err, malak.ErrContactNotFoundOnBoard)
		require.Nil(t, result)
	})
}

func TestFundraising_UpdateContactDeal(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	fundingRepo := NewFundingRepo(client)
	workspaceRepo := NewWorkspaceRepository(client)

	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	pipeline := &malak.FundraisingPipeline{
		Reference:         malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipeline),
		WorkspaceID:       workspace.ID,
		Title:             "Test Pipeline",
		Stage:             malak.FundraisePipelineStageSeed,
		Description:       "Test pipeline description",
		TargetAmount:      1000000,
		StartDate:         time.Now().UTC(),
		ExpectedCloseDate: time.Now().UTC().Add(90 * 24 * time.Hour),
	}

	column := malak.FundraisingPipelineColumn{
		Title:       "Test Column",
		ColumnType:  malak.FundraisePipelineColumnTypeNormal,
		Description: "Test column description",
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeFundraisingPipelineColumn),
	}

	err = fundingRepo.Create(t.Context(), pipeline, column)
	require.NoError(t, err)

	contact := &malak.Contact{
		ID:          uuid.New(),
		Email:       malak.Email("test@example.com"),
		WorkspaceID: workspace.ID,
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeContact),
		FirstName:   "Test",
		LastName:    "Contact",
	}

	_, err = client.NewInsert().Model(contact).Exec(t.Context())
	require.NoError(t, err)

	var columns []malak.FundraisingPipelineColumn
	err = client.NewSelect().
		Model(&columns).
		Where("fundraising_pipeline_id = ?", pipeline.ID).
		Scan(t.Context())
	require.NoError(t, err)
	require.Len(t, columns, 1)

	// Add contact to board first
	opts := &malak.AddContactToBoardOptions{
		Column:             &columns[0],
		Contact:            contact,
		ReferenceGenerator: malak.NewReferenceGenerator(),
		Rating:             3,
		CanLeadRound:       true,
		CheckSize:          1000000, // $1M
	}

	err = fundingRepo.AddContactToBoard(t.Context(), opts)
	require.NoError(t, err)

	// Get the fundraise contact ID
	var fundraiseContact malak.FundraiseContact
	err = client.NewSelect().
		Model(&fundraiseContact).
		Where("contact_id = ?", contact.ID).
		Where("fundraising_pipeline_id = ?", pipeline.ID).
		Scan(t.Context())
	require.NoError(t, err)

	t.Run("update contact deal details", func(t *testing.T) {
		updateOpts := malak.UpdateContactDealOptions{
			ContactID:    fundraiseContact.ID,
			Rating:       5,
			CanLeadRound: false,
			CheckSize:    5000000, // $5M
		}

		err := fundingRepo.UpdateContactDeal(t.Context(), pipeline, updateOpts)
		require.NoError(t, err)

		// Verify the update
		var dealDetails malak.FundraiseContactDealDetails
		err = client.NewSelect().
			Model(&dealDetails).
			Where("fundraising_pipeline_column_contact_id = ?", fundraiseContact.ID).
			Scan(t.Context())
		require.NoError(t, err)
		require.Equal(t, int64(5), dealDetails.Rating)
		require.Equal(t, false, dealDetails.CanLeadRound)
		require.Equal(t, int64(5000000), dealDetails.CheckSize)
	})

	t.Run("update non-existent contact deal", func(t *testing.T) {
		updateOpts := malak.UpdateContactDealOptions{
			ContactID:    uuid.New(),
			Rating:       4,
			CanLeadRound: true,
			CheckSize:    1000000, // $1M
		}

		err := fundingRepo.UpdateContactDeal(t.Context(), pipeline, updateOpts)
		require.NoError(t, err) // Should not error as it's a no-op update
	})
}
