// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package malak

import (
	"errors"
	"fmt"
)

const (
	// ListUpdateFilterStatusDraft is a ListUpdateFilterStatus of type draft.
	ListUpdateFilterStatusDraft ListUpdateFilterStatus = "draft"
	// ListUpdateFilterStatusSent is a ListUpdateFilterStatus of type sent.
	ListUpdateFilterStatusSent ListUpdateFilterStatus = "sent"
	// ListUpdateFilterStatusAll is a ListUpdateFilterStatus of type all.
	ListUpdateFilterStatusAll ListUpdateFilterStatus = "all"
)

var ErrInvalidListUpdateFilterStatus = errors.New("not a valid ListUpdateFilterStatus")

// String implements the Stringer interface.
func (x ListUpdateFilterStatus) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x ListUpdateFilterStatus) IsValid() bool {
	_, err := ParseListUpdateFilterStatus(string(x))
	return err == nil
}

var _ListUpdateFilterStatusValue = map[string]ListUpdateFilterStatus{
	"draft": ListUpdateFilterStatusDraft,
	"sent":  ListUpdateFilterStatusSent,
	"all":   ListUpdateFilterStatusAll,
}

// ParseListUpdateFilterStatus attempts to convert a string to a ListUpdateFilterStatus.
func ParseListUpdateFilterStatus(name string) (ListUpdateFilterStatus, error) {
	if x, ok := _ListUpdateFilterStatusValue[name]; ok {
		return x, nil
	}
	return ListUpdateFilterStatus(""), fmt.Errorf("%s is %w", name, ErrInvalidListUpdateFilterStatus)
}

const (
	// ReactionStatusThumbsUp is a ReactionStatus of type thumbs up.
	ReactionStatusThumbsUp ReactionStatus = "thumbs up"
	// ReactionStatusThumbsDown is a ReactionStatus of type thumbs down.
	ReactionStatusThumbsDown ReactionStatus = "thumbs down"
)

var ErrInvalidReactionStatus = errors.New("not a valid ReactionStatus")

// String implements the Stringer interface.
func (x ReactionStatus) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x ReactionStatus) IsValid() bool {
	_, err := ParseReactionStatus(string(x))
	return err == nil
}

var _ReactionStatusValue = map[string]ReactionStatus{
	"thumbs up":   ReactionStatusThumbsUp,
	"thumbs down": ReactionStatusThumbsDown,
}

// ParseReactionStatus attempts to convert a string to a ReactionStatus.
func ParseReactionStatus(name string) (ReactionStatus, error) {
	if x, ok := _ReactionStatusValue[name]; ok {
		return x, nil
	}
	return ReactionStatus(""), fmt.Errorf("%s is %w", name, ErrInvalidReactionStatus)
}

const (
	// RecipientStatusPending is a RecipientStatus of type pending.
	RecipientStatusPending RecipientStatus = "pending"
	// RecipientStatusSent is a RecipientStatus of type sent.
	RecipientStatusSent RecipientStatus = "sent"
	// RecipientStatusFailed is a RecipientStatus of type failed.
	RecipientStatusFailed RecipientStatus = "failed"
)

var ErrInvalidRecipientStatus = errors.New("not a valid RecipientStatus")

// String implements the Stringer interface.
func (x RecipientStatus) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x RecipientStatus) IsValid() bool {
	_, err := ParseRecipientStatus(string(x))
	return err == nil
}

var _RecipientStatusValue = map[string]RecipientStatus{
	"pending": RecipientStatusPending,
	"sent":    RecipientStatusSent,
	"failed":  RecipientStatusFailed,
}

// ParseRecipientStatus attempts to convert a string to a RecipientStatus.
func ParseRecipientStatus(name string) (RecipientStatus, error) {
	if x, ok := _RecipientStatusValue[name]; ok {
		return x, nil
	}
	return RecipientStatus(""), fmt.Errorf("%s is %w", name, ErrInvalidRecipientStatus)
}

const (
	// RecipientTypeList is a RecipientType of type list.
	RecipientTypeList RecipientType = "list"
	// RecipientTypeEmail is a RecipientType of type email.
	RecipientTypeEmail RecipientType = "email"
)

var ErrInvalidRecipientType = errors.New("not a valid RecipientType")

// String implements the Stringer interface.
func (x RecipientType) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x RecipientType) IsValid() bool {
	_, err := ParseRecipientType(string(x))
	return err == nil
}

var _RecipientTypeValue = map[string]RecipientType{
	"list":  RecipientTypeList,
	"email": RecipientTypeEmail,
}

// ParseRecipientType attempts to convert a string to a RecipientType.
func ParseRecipientType(name string) (RecipientType, error) {
	if x, ok := _RecipientTypeValue[name]; ok {
		return x, nil
	}
	return RecipientType(""), fmt.Errorf("%s is %w", name, ErrInvalidRecipientType)
}

const (
	// UpdateRecipientLogProviderResend is a UpdateRecipientLogProvider of type resend.
	UpdateRecipientLogProviderResend UpdateRecipientLogProvider = "resend"
	// UpdateRecipientLogProviderSendgrid is a UpdateRecipientLogProvider of type sendgrid.
	UpdateRecipientLogProviderSendgrid UpdateRecipientLogProvider = "sendgrid"
	// UpdateRecipientLogProviderSmtp is a UpdateRecipientLogProvider of type smtp.
	UpdateRecipientLogProviderSmtp UpdateRecipientLogProvider = "smtp"
)

var ErrInvalidUpdateRecipientLogProvider = errors.New("not a valid UpdateRecipientLogProvider")

// String implements the Stringer interface.
func (x UpdateRecipientLogProvider) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x UpdateRecipientLogProvider) IsValid() bool {
	_, err := ParseUpdateRecipientLogProvider(string(x))
	return err == nil
}

var _UpdateRecipientLogProviderValue = map[string]UpdateRecipientLogProvider{
	"resend":   UpdateRecipientLogProviderResend,
	"sendgrid": UpdateRecipientLogProviderSendgrid,
	"smtp":     UpdateRecipientLogProviderSmtp,
}

// ParseUpdateRecipientLogProvider attempts to convert a string to a UpdateRecipientLogProvider.
func ParseUpdateRecipientLogProvider(name string) (UpdateRecipientLogProvider, error) {
	if x, ok := _UpdateRecipientLogProviderValue[name]; ok {
		return x, nil
	}
	return UpdateRecipientLogProvider(""), fmt.Errorf("%s is %w", name, ErrInvalidUpdateRecipientLogProvider)
}

const (
	// UpdateSendScheduleScheduled is a UpdateSendSchedule of type scheduled.
	UpdateSendScheduleScheduled UpdateSendSchedule = "scheduled"
	// UpdateSendScheduleCancelled is a UpdateSendSchedule of type cancelled.
	UpdateSendScheduleCancelled UpdateSendSchedule = "cancelled"
	// UpdateSendScheduleSent is a UpdateSendSchedule of type sent.
	UpdateSendScheduleSent UpdateSendSchedule = "sent"
	// UpdateSendScheduleFailed is a UpdateSendSchedule of type failed.
	UpdateSendScheduleFailed UpdateSendSchedule = "failed"
	// UpdateSendScheduleProcessing is a UpdateSendSchedule of type processing.
	UpdateSendScheduleProcessing UpdateSendSchedule = "processing"
)

var ErrInvalidUpdateSendSchedule = errors.New("not a valid UpdateSendSchedule")

// String implements the Stringer interface.
func (x UpdateSendSchedule) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x UpdateSendSchedule) IsValid() bool {
	_, err := ParseUpdateSendSchedule(string(x))
	return err == nil
}

var _UpdateSendScheduleValue = map[string]UpdateSendSchedule{
	"scheduled":  UpdateSendScheduleScheduled,
	"cancelled":  UpdateSendScheduleCancelled,
	"sent":       UpdateSendScheduleSent,
	"failed":     UpdateSendScheduleFailed,
	"processing": UpdateSendScheduleProcessing,
}

// ParseUpdateSendSchedule attempts to convert a string to a UpdateSendSchedule.
func ParseUpdateSendSchedule(name string) (UpdateSendSchedule, error) {
	if x, ok := _UpdateSendScheduleValue[name]; ok {
		return x, nil
	}
	return UpdateSendSchedule(""), fmt.Errorf("%s is %w", name, ErrInvalidUpdateSendSchedule)
}

const (
	// UpdateStatusDraft is a UpdateStatus of type draft.
	UpdateStatusDraft UpdateStatus = "draft"
	// UpdateStatusSent is a UpdateStatus of type sent.
	UpdateStatusSent UpdateStatus = "sent"
)

var ErrInvalidUpdateStatus = errors.New("not a valid UpdateStatus")

// String implements the Stringer interface.
func (x UpdateStatus) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x UpdateStatus) IsValid() bool {
	_, err := ParseUpdateStatus(string(x))
	return err == nil
}

var _UpdateStatusValue = map[string]UpdateStatus{
	"draft": UpdateStatusDraft,
	"sent":  UpdateStatusSent,
}

// ParseUpdateStatus attempts to convert a string to a UpdateStatus.
func ParseUpdateStatus(name string) (UpdateStatus, error) {
	if x, ok := _UpdateStatusValue[name]; ok {
		return x, nil
	}
	return UpdateStatus(""), fmt.Errorf("%s is %w", name, ErrInvalidUpdateStatus)
}

const (
	// UpdateTypePreview is a UpdateType of type preview.
	UpdateTypePreview UpdateType = "preview"
	// UpdateTypeLive is a UpdateType of type live.
	UpdateTypeLive UpdateType = "live"
)

var ErrInvalidUpdateType = errors.New("not a valid UpdateType")

// String implements the Stringer interface.
func (x UpdateType) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x UpdateType) IsValid() bool {
	_, err := ParseUpdateType(string(x))
	return err == nil
}

var _UpdateTypeValue = map[string]UpdateType{
	"preview": UpdateTypePreview,
	"live":    UpdateTypeLive,
}

// ParseUpdateType attempts to convert a string to a UpdateType.
func ParseUpdateType(name string) (UpdateType, error) {
	if x, ok := _UpdateTypeValue[name]; ok {
		return x, nil
	}
	return UpdateType(""), fmt.Errorf("%s is %w", name, ErrInvalidUpdateType)
}
