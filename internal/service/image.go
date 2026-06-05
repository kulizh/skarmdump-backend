package service

import (
	"bytes"

	"skarmdump-backend/pkg/hash"
)

type Local interface {
	Save(hash string, data []byte) error
}

type S3 interface {
	Save(hash string, data []byte) error
}

type Service struct {
	local Local
	s3    S3
}

func New(l Local, s3 S3) *Service {
	return &Service{local: l, s3: s3}
}

func (s *Service) Upload(data []byte, useS3 bool) (string, error) {
	h, _, err := hash.ShortHash(bytes.NewReader(data), 8)
	if err != nil {
		return "", err
	}

	if useS3 {
		if err := s.s3.Save(h, data); err != nil {
			return "", err
		}
	} else {
		if err := s.local.Save(h, data); err != nil {
			return "", err
		}
	}

	return h, nil
}
