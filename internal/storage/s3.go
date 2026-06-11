package storage

import (
	"bytes"
	"context"
	"fmt"
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

func s3key(hash, subdir, ext string) string {
	if subdir == "" {
		return fmt.Sprintf("%s.%s", hash, ext)
	}
	return fmt.Sprintf("%s/%s.%s", subdir, hash, ext)
}

func (s *S3) SavePNG(hash string, data []byte) error {
	return s.put(hash, "", "png", data)
}

func (s *S3) SaveOG(hash string, data []byte) error {
	return s.put(hash, "og", "png", data)
}

func (s *S3) SaveRE(hash string, data []byte) error {
	return s.put(hash, "re", "png", data)
}

func (s *S3) put(hash, subdir, ext string, data []byte) error {
	key := s3key(hash, subdir, ext)
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
	return s.get(hash, "", "png")
}

func (s *S3) ReadOG(hash string) ([]byte, error) {
	return s.get(hash, "og", "png")
}

func (s *S3) ReadRE(hash string) ([]byte, error) {
	return s.get(hash, "re", "png")
}

func (s *S3) get(hash, subdir, ext string) ([]byte, error) {
	key := s3key(hash, subdir, ext)
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
