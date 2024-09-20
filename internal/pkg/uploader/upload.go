package uploader

import (
	"context"
	"io"
)

type UploadedFileMetadata struct {
	FolderDestination string `json:"folder_destination,omitempty"`
	Key               string `json:"key,omitempty"`
	Size              int64  `json:"size,omitempty"`
}

type Uploader interface {
	Upload(context.Context, io.Reader) (*UploadedFileMetadata, error)
	io.Closer
}
