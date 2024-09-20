package uploader

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/google/uuid"
)

type S3Options struct {
	Bucket string
}

type S3Store struct {
	client *s3.Client
	opts   S3Options
}

func NewS3FromConfig(cfg aws.Config, opts S3Options) (*S3Store, error) {

	if util.IsStringEmpty(opts.Bucket) {
		return nil, errors.New("please provide a valid s3 bucket")
	}

	return &S3Store{
		client: s3.NewFromConfig(cfg),
	}, nil
}

func (s *S3Store) Close() error { return nil }

func (s *S3Store) Upload(ctx context.Context, r io.Reader,
) (*UploadedFileMetadata, error) {
	b := new(bytes.Buffer)

	r = io.TeeReader(r, b)

	n, err := io.Copy(io.Discard, r)
	if err != nil {
		return nil, err
	}

	key := uuid.NewString()

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: util.Ref(s.opts.Bucket),
		Key:    util.Ref(key),
		Body:   b,
	})
	if err != nil {
		return nil, err
	}

	return &UploadedFileMetadata{
		FolderDestination: s.opts.Bucket,
		Size:              n,
		Key:               key,
	}, nil
}
