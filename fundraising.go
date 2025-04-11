package malak

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var (
	ErrPipelineNotFound       = errors.New("pipeline not found")
	ErrContactNotFoundOnBoard = errors.New("contact not found on board")
	ErrPipelineColumnNotFound = errors.New("column not found in pipeline")
)

var DefaultFundraisingColumns = []struct {
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

// ENUM(family_and_friend,pre_seed,bridge_round,seed,series_a,series_b,series_c)
type FundraisePipelineStage string

// ENUM(active,closed)
type FundraisePipelineStatus string

// ENUM(normal,closed)
//
// normal columns are just normal columns really
// but closed columns are those that signify the investor in the pipeline is now "closed"
// that might mean documents signed, money wired. It depends on you
// But there can only be one closed column per fundraising pipeline
type FundraisePipelineColumnType string

type FundraisingPipeline struct {
	ID uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`

	Stage        FundraisePipelineStage `json:"stage,omitempty"`
	Reference    Reference              `json:"reference,omitempty"`
	WorkspaceID  uuid.UUID              `json:"workspace_id,omitempty"`
	Title        string                 `json:"title,omitempty"`
	Description  string                 `json:"description,omitempty"`
	IsClosed     bool                   `json:"is_closed,omitempty"`
	TargetAmount int64                  `json:"target_amount,omitempty"`
	// this is being updated dynamically by postgres triggers
	// We also use to calculate progress
	ClosedAmount int64 `json:"closed_amount,omitempty"`

	// Can be in the future
	StartDate         time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"start_date,omitempty"`
	ExpectedCloseDate time.Time `json:"expected_close_date,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`

	bun.BaseModel `json:"-"`
}

type FundraisingPipelineColumn struct {
	ID uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`

	Reference             Reference                   `json:"reference,omitempty"`
	FundraisingPipelineID uuid.UUID                   `json:"fundraising_pipeline_id,omitempty"`
	Title                 string                      `json:"title,omitempty"`
	ColumnType            FundraisePipelineColumnType `json:"column_type,omitempty"`
	Description           string                      `json:"description,omitempty"`
	InvestorsCount        int64                       `json:"investors_count,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`

	bun.BaseModel `json:"-"`
}

type FundraiseContactPosition struct {
	ID                                 uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Reference                          Reference `json:"reference,omitempty"`
	FundraisingPipelineColumnContactID uuid.UUID `json:"fundraising_pipeline_column_contact_id,omitempty"`
	OrderIndex                         int64     `json:"order_index,omitempty"`

	bun.BaseModel `json:"-" bun:"table:fundraising_pipeline_column_contact_positions"`
}

type FundraiseContact struct {
	ID                          uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Reference                   Reference `json:"reference,omitempty"`
	ContactID                   uuid.UUID `json:"contact_id,omitempty"`
	FundraisingPipelineID       uuid.UUID `json:"fundraising_pipeline_id,omitempty"`
	FundraisingPipelineColumnID uuid.UUID `json:"fundraising_pipeline_column_id,omitempty"`

	Contact     *Contact                     `bun:"rel:belongs-to,join:contact_id=id" json:"contact,omitempty"`
	DealDetails *FundraiseContactDealDetails `bun:"rel:has-one,join:id=fundraising_pipeline_column_contact_id" json:"deal_details,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`

	bun.BaseModel `json:"-" bun:"table:fundraising_pipeline_column_contacts"`
}

type FundraiseContactDealDetails struct {
	ID                                 uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Reference                          Reference `json:"reference,omitempty"`
	FundraisingPipelineColumnContactID uuid.UUID `json:"fundraising_pipeline_column_contact_id,omitempty"`
	CheckSize                          int64     `json:"check_size,omitempty"`
	CanLeadRound                       bool      `json:"can_lead_round,omitempty"`
	Rating                             int64     `json:"rating,omitempty"`
	InitialContact                     time.Time `json:"initial_contact,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`

	bun.BaseModel `json:"-" bun:"table:fundraising_pipeline_column_contact_deals"`
}

// ENUM(meeting,note,email)
type FundraisingColumnActivity string

type FundraiseContactActivity struct {
	ID                                 uuid.UUID                 `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Reference                          Reference                 `json:"reference,omitempty"`
	FundraisingPipelineColumnContactID uuid.UUID                 `json:"fundraising_pipeline_column_contact_id,omitempty"`
	ActivityType                       FundraisingColumnActivity `json:"activity_type,omitempty"`
	Title                              string                    `json:"title,omitempty"`
	Content                            string                    `json:"content,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`

	bun.BaseModel `json:"-" bun:"table:fundraising_pipeline_column_contact_activities"`
}

type FundraiseContactDocument struct {
	ID                                 uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Reference                          Reference `json:"reference,omitempty"`
	FundraisingPipelineColumnContactID uuid.UUID `json:"fundraising_pipeline_column_contact_id,omitempty"`
	Title                              string    `json:"title,omitempty"`
	FileSize                           int64     `json:"file_size,omitempty"`
	ObjectKey                          string    `json:"object_key,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`

	bun.BaseModel `json:"-" bun:"table:fundraising_pipeline_column_contact_documents"`
}

type ListPipelineOptions struct {
	Paginator   Paginator
	WorkspaceID uuid.UUID
	ActiveOnly  bool
}

type FetchPipelineOptions struct {
	WorkspaceID uuid.UUID
	Reference   Reference
}

type AddContactToBoardOptions struct {
	Column             *FundraisingPipelineColumn
	Contact            *Contact
	ReferenceGenerator ReferenceGeneratorOperation
	Rating             int
	CanLeadRound       bool
	InitialContact     time.Time
	CheckSize          int64
}

type UpdateContactDealOptions struct {
	Rating       int64
	CanLeadRound bool
	CheckSize    int64
	ContactID    uuid.UUID
}

type GetBoardOptions struct {
	PipelineID uuid.UUID
	ColumnID   uuid.UUID
}

type FundraisingPipelineRepository interface {
	Create(context.Context, *FundraisingPipeline, ...FundraisingPipelineColumn) error
	List(context.Context, ListPipelineOptions) ([]FundraisingPipeline, int64, error)
	Get(context.Context, FetchPipelineOptions) (*FundraisingPipeline, error)
	Board(context.Context, *FundraisingPipeline) ([]FundraisingPipelineColumn, []FundraiseContact, []FundraiseContactPosition, error)
	CloseBoard(context.Context, *FundraisingPipeline) error

	// This is just the first inserted column for now. keeping it simple
	DefaultColumn(context.Context, *FundraisingPipeline) (FundraisingPipelineColumn, error)

	GetColumn(context.Context, GetBoardOptions) (*FundraisingPipelineColumn, error)

	AddContactToBoard(context.Context, *AddContactToBoardOptions) error
	UpdateBoardContact(context.Context, *FundraiseContactDealDetails) error
	GetContact(context.Context, uuid.UUID, uuid.UUID) (*FundraiseContact, error)
	UpdateContactDeal(context.Context, *FundraisingPipeline, UpdateContactDealOptions) error
}
