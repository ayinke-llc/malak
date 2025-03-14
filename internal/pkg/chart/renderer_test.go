package chart

import (
	"testing"

	"github.com/adelowo/gulter"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	malak_mocks "github.com/ayinke-llc/malak/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"go.uber.org/mock/gomock"
)

func TestEChartsRenderer_RenderChart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := NewMockStorage(ctrl)
	mockIntegration := malak_mocks.NewMockIntegrationRepository(ctrl)
	db := &bun.DB{}
	cfg := config.Config{
		Uploader: struct {
			Driver        config.UploadDriver `yaml:"driver" mapstructure:"driver"`
			MaxUploadSize int64               `yaml:"max_upload_size" mapstructure:"max_upload_size"`
			S3            struct {
				AccessKey     string `yaml:"access_key" mapstructure:"access_key"`
				AccessSecret  string `yaml:"access_secret" mapstructure:"access_secret"`
				Region        string `yaml:"region" mapstructure:"region"`
				Endpoint      string `yaml:"endpoint" mapstructure:"endpoint"`
				LogOperations bool   `yaml:"log_operations" mapstructure:"log_operations"`
				Bucket        string `yaml:"bucket" mapstructure:"bucket"`
				DeckBucket    string `yaml:"deck_bucket" mapstructure:"deck_bucket"`
				UseTLS        bool   `yaml:"use_tls" mapstructure:"use_tls"`
			} `yaml:"s3" mapstructure:"s3"`
		}{
			S3: struct {
				AccessKey     string `yaml:"access_key" mapstructure:"access_key"`
				AccessSecret  string `yaml:"access_secret" mapstructure:"access_secret"`
				Region        string `yaml:"region" mapstructure:"region"`
				Endpoint      string `yaml:"endpoint" mapstructure:"endpoint"`
				LogOperations bool   `yaml:"log_operations" mapstructure:"log_operations"`
				Bucket        string `yaml:"bucket" mapstructure:"bucket"`
				DeckBucket    string `yaml:"deck_bucket" mapstructure:"deck_bucket"`
				UseTLS        bool   `yaml:"use_tls" mapstructure:"use_tls"`
			}{
				Bucket:   "test-bucket",
				Endpoint: "http://test-endpoint",
			},
		},
	}

	tests := []struct {
		name          string
		workspaceID   uuid.UUID
		chartID       string
		setupMocks    func(storage *MockStorage, integration *malak_mocks.MockIntegrationRepository)
		expectedError bool
		expectedURL   string
	}{
		{
			name:        "chart not found",
			workspaceID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			chartID:     "chart123",
			setupMocks: func(storage *MockStorage, integration *malak_mocks.MockIntegrationRepository) {
				integration.EXPECT().
					GetChart(gomock.Any(), malak.FetchChartOptions{
						WorkspaceID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
						Reference:   malak.Reference("chart123"),
					}).
					Return(malak.IntegrationChart{}, malak.ErrChartNotFound)
			},
			expectedError: true,
		},
		{
			name:        "unsupported chart type",
			workspaceID: uuid.New(),
			chartID:     "chart123",
			setupMocks: func(storage *MockStorage, integration *malak_mocks.MockIntegrationRepository) {
				integration.EXPECT().
					GetChart(gomock.Any(), gomock.Any()).
					Return(malak.IntegrationChart{
						ChartType: "unsupported",
					}, nil)
			},
			expectedError: true,
		},
		{
			name:        "successful bar chart with currency values",
			workspaceID: uuid.New(),
			chartID:     "chart123",
			setupMocks: func(storage *MockStorage, integration *malak_mocks.MockIntegrationRepository) {
				chart := malak.IntegrationChart{
					ChartType: "bar",
				}

				integration.EXPECT().
					GetChart(gomock.Any(), gomock.Any()).
					Return(chart, nil)

				integration.EXPECT().
					GetDataPoints(gomock.Any(), gomock.Any()).
					Return([]malak.IntegrationDataPoint{
						{
							PointName:     "Revenue",
							PointValue:    2540,
							DataPointType: "currency",
						},
						{
							PointName:     "Expenses",
							PointValue:    1250,
							DataPointType: "currency",
						},
					}, nil)

				storage.EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&gulter.UploadedFileMetadata{
						FolderDestination: "charts",
						Key:               "test-key.png",
					}, nil)
			},
			expectedError: false,
			expectedURL:   "http://test-endpoint/charts/test-key.png",
		},
		{
			name:        "successful bar chart with mixed data types",
			workspaceID: uuid.New(),
			chartID:     "chart123",
			setupMocks: func(storage *MockStorage, integration *malak_mocks.MockIntegrationRepository) {
				chart := malak.IntegrationChart{
					ChartType: "bar",
				}

				integration.EXPECT().
					GetChart(gomock.Any(), gomock.Any()).
					Return(chart, nil)

				integration.EXPECT().
					GetDataPoints(gomock.Any(), gomock.Any()).
					Return([]malak.IntegrationDataPoint{
						{
							PointName:     "Revenue",
							PointValue:    2540,
							DataPointType: "currency",
						},
						{
							PointName:     "Count",
							PointValue:    50,
							DataPointType: "others",
						},
					}, nil)

				storage.EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&gulter.UploadedFileMetadata{
						FolderDestination: "charts",
						Key:               "test-key.png",
					}, nil)
			},
			expectedError: false,
			expectedURL:   "http://test-endpoint/charts/test-key.png",
		},
		{
			name:        "successful pie chart with currency values",
			workspaceID: uuid.New(),
			chartID:     "chart123",
			setupMocks: func(storage *MockStorage, integration *malak_mocks.MockIntegrationRepository) {
				chart := malak.IntegrationChart{
					ChartType: "pie",
				}

				integration.EXPECT().
					GetChart(gomock.Any(), gomock.Any()).
					Return(chart, nil)

				integration.EXPECT().
					GetDataPoints(gomock.Any(), gomock.Any()).
					Return([]malak.IntegrationDataPoint{
						{
							PointName:     "Revenue",
							PointValue:    2540,
							DataPointType: "currency",
						},
						{
							PointName:     "Expenses",
							PointValue:    1250,
							DataPointType: "currency",
						},
					}, nil)

				storage.EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&gulter.UploadedFileMetadata{
						FolderDestination: "charts",
						Key:               "test-key.png",
					}, nil)
			},
			expectedError: false,
			expectedURL:   "http://test-endpoint/charts/test-key.png",
		},
		{
			name:        "successful pie chart with mixed data types",
			workspaceID: uuid.New(),
			chartID:     "chart123",
			setupMocks: func(storage *MockStorage, integration *malak_mocks.MockIntegrationRepository) {
				chart := malak.IntegrationChart{
					ChartType: "pie",
				}

				integration.EXPECT().
					GetChart(gomock.Any(), gomock.Any()).
					Return(chart, nil)

				integration.EXPECT().
					GetDataPoints(gomock.Any(), gomock.Any()).
					Return([]malak.IntegrationDataPoint{
						{
							PointName:     "Revenue",
							PointValue:    2540,
							DataPointType: "currency",
						},
						{
							PointName:     "Count",
							PointValue:    50,
							DataPointType: "others",
						},
					}, nil)

				storage.EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&gulter.UploadedFileMetadata{
						FolderDestination: "charts",
						Key:               "test-key.png",
					}, nil)
			},
			expectedError: false,
			expectedURL:   "http://test-endpoint/charts/test-key.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewEChartsRenderer(mockStorage, cfg, db, mockIntegration)
			tt.setupMocks(mockStorage, mockIntegration)

			url, err := renderer.RenderChart(tt.workspaceID, tt.chartID)
			if tt.expectedError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedURL, url)
		})
	}
}

func TestEChartsRenderer_renderBarChart(t *testing.T) {
	renderer := NewEChartsRenderer(nil, config.Config{}, nil, nil)

	tests := []struct {
		name       string
		dataPoints []ChartDataPoint
		wantErr    bool
		validateFn func(t *testing.T, output []byte)
	}{
		{
			name: "renders bar chart with currency values",
			dataPoints: []ChartDataPoint{
				{Label: "Revenue", Value: 25.40},
				{Label: "Expenses", Value: 12.50},
			},
			wantErr: false,
			validateFn: func(t *testing.T, output []byte) {
				require.Contains(t, string(output), "Revenue")
				require.Contains(t, string(output), "Expenses")
				require.Contains(t, string(output), "25.4")
				require.Contains(t, string(output), "12.5")
			},
		},
		{
			name: "renders bar chart with mixed values",
			dataPoints: []ChartDataPoint{
				{Label: "Revenue", Value: 25.40},
				{Label: "Count", Value: 50},
			},
			wantErr: false,
			validateFn: func(t *testing.T, output []byte) {
				require.Contains(t, string(output), "Revenue")
				require.Contains(t, string(output), "Count")
				require.Contains(t, string(output), "25.4")
				require.Contains(t, string(output), "50")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := renderer.renderBarChart(tt.dataPoints)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			tt.validateFn(t, output)
		})
	}
}

func TestEChartsRenderer_renderPieChart(t *testing.T) {
	renderer := NewEChartsRenderer(nil, config.Config{}, nil, nil)

	tests := []struct {
		name       string
		dataPoints []ChartDataPoint
		wantErr    bool
		validateFn func(t *testing.T, output []byte)
	}{
		{
			name: "renders pie chart with currency values",
			dataPoints: []ChartDataPoint{
				{Label: "Revenue", Value: 25.40},
				{Label: "Expenses", Value: 12.50},
			},
			wantErr: false,
			validateFn: func(t *testing.T, output []byte) {
				require.Contains(t, string(output), "Revenue")
				require.Contains(t, string(output), "Expenses")
				require.Contains(t, string(output), "25.4")
				require.Contains(t, string(output), "12.5")
			},
		},
		{
			name: "renders pie chart with mixed values",
			dataPoints: []ChartDataPoint{
				{Label: "Revenue", Value: 25.40},
				{Label: "Count", Value: 50},
			},
			wantErr: false,
			validateFn: func(t *testing.T, output []byte) {
				require.Contains(t, string(output), "Revenue")
				require.Contains(t, string(output), "Count")
				require.Contains(t, string(output), "25.4")
				require.Contains(t, string(output), "50")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := renderer.renderPieChart(tt.dataPoints)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			tt.validateFn(t, output)
		})
	}
}
