package malak

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var defaultFundraisingColumns = []struct {
	Title       string
	ColumnType  FundraisePipelineColumnType
	Description string
}{
	{
		Title:       "Backlog",
		ColumnType:  FundraisePipelineColumnTypeNormal,
		Description: "Investors you would love to speak to",
	},
	{
		Title:       "Contacted",
		ColumnType:  FundraisePipelineColumnTypeNormal,
		Description: "Investors you contacted or submitted application via their website or those that reached out to you",
	},
	{
		Title:       "Partner Meeting",
		ColumnType:  FundraisePipelineColumnTypeNormal,
		Description: "Investors you are currently speaking to",
	},
	{
		Title:       "Passed",
		ColumnType:  FundraisePipelineColumnTypeNormal,
		Description: "Investors you spoke to but it didn't pane out",
	},
	{
		Title:       "Termsheet/SAFE",
		ColumnType:  FundraisePipelineColumnTypeNormal,
		Description: "Investors that have given you a termsheet or SAFE",
	},
	{
		Title:       "Closed",
		ColumnType:  FundraisePipelineColumnTypeClosed,
		Description: "Investors that have signed the termsheet/safe and closed the deal. This might mean money wired or not",
	},
}

// ENUM(family_and_friend,pre_seed,seed,series_a,series_b,series_c)
type FundraisePipelineStage string

// ENUM(normal,closed)
//
// normal columns are just normal columns really
// but closed columns are those that signify the investor in the pipeline is now "closed"
// that might mean documents signed, money wired. It depends on you
// But there can only be one closed column per fundraising pipeline
type FundraisePipelineColumnType uint8

type FundraisingPipeline struct {
	ID uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`

	Reference        Reference `json:"reference,omitempty"`
	WorkspaceID      uuid.UUID `json:"workspace_id,omitempty"`
	CreatedBy        uuid.UUID `json:"created_by,omitempty"`
	Title            string    `json:"title,omitempty"`
	Description      string    `json:"description,omitempty"`
	ExpectedDeadline time.Time `json:"expected_deadline,omitempty"`

	TargetAmount int64 `json:"target_amount,omitempty"`

	// this is being updated dynamically by postgres triggers
	// We also use to calculate progress
	ClosedAmount int64 `json:"closed_amount,omitempty"`

	IsClosed bool `json:"is_closed,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`

	bun.BaseModel `bun:"table:decks" json:"-"`
}

type FundraisingPipelineColumn struct {
	ID uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`

	Reference             Reference                   `json:"reference,omitempty"`
	FundraisingPipelineID uuid.UUID                   `json:"fundraising_pipeline_id,omitempty"`
	Title                 string                      `json:"title,omitempty"`
	ColumnType            FundraisePipelineColumnType `json:"column_type,omitempty"`
	Description           string                      `json:"description,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`

	bun.BaseModel `bun:"table:decks" json:"-"`
}

type FundraisingPipelineRepository interface{}
