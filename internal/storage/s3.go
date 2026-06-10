package storage

import (
	"bytes"
	"context"
	"io"

	"skarmdump-backend/internal/config"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	client *s3.Client
	cfg    *config.Config
}

func NewS3(cfg *config.Config) *S3 {
	return &S3{
		client: s3.New(s3.Options{
			Region: cfg.S3Region,
		}),
		cfg: cfg,
	}
}

func (s *S3) Save(hash string, data []byte) error {
	key := hash + ".png"

	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s.cfg.S3Bucket,
		Key:    &key,
		Body:   bytes.NewReader(data),
	})
	return err
}

func (s *S3) Exists(hash string) bool {
	key := hash + ".png"
	_, err := s.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: &s.cfg.S3Bucket,
		Key:    &key,
	})
	return err == nil
}

func (s *S3) Read(hash string) ([]byte, error) {
	key := hash + ".png"
	out, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &s.cfg.S3Bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer out.Body.Close()
	return io.ReadAll(out.Body)
}
