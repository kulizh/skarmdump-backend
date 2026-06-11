package service

import (
	"bytes"

	"skarmdump-backend/pkg/hash"
)

type Local interface {
	SavePNG(hash string, data []byte) error
	SaveOG(hash string, data []byte) error
	SaveRE(hash string, data []byte) error
	Exists(hash string) bool
	Read(hash string) ([]byte, error)
	ReadOG(hash string) ([]byte, error)
	ReadRE(hash string) ([]byte, error)
}

type S3 interface {
	SavePNG(hash string, data []byte) error
	SaveOG(hash string, data []byte) error
	SaveRE(hash string, data []byte) error
	Exists(hash string) bool
	Read(hash string) ([]byte, error)
	ReadOG(hash string) ([]byte, error)
	ReadRE(hash string) ([]byte, error)
}

type Service struct {
	local      Local
	s3         S3
	hashLength int
	resizer    *Resizer
}

func New(l Local, s3 S3, hashLength int, resizer *Resizer) *Service {
	return &Service{local: l, s3: s3, hashLength: hashLength, resizer: resizer}
}

func (s *Service) Exists(hash string) bool {
	return s.local.Exists(hash) || s.s3.Exists(hash)
}

func (s *Service) Get(hash string) ([]byte, error) {
	if s.local.Exists(hash) {
		return s.local.Read(hash)
	}
	return s.s3.Read(hash)
}

func (s *Service) GetOG(hash string) ([]byte, error) {
	if s.local.Exists(hash) {
		return s.local.ReadOG(hash)
	}
	return s.s3.ReadOG(hash)
}

func (s *Service) GetRE(hash string) ([]byte, error) {
	if s.local.Exists(hash) {
		return s.local.ReadRE(hash)
	}
	return s.s3.ReadRE(hash)
}

func (s *Service) Upload(data []byte, useS3 bool) (string, error) {
	h, _, err := hash.ShortHash(bytes.NewReader(data), s.hashLength)
	if err != nil {
		return "", err
	}

	if useS3 {
		if err := s.s3.SavePNG(h, data); err != nil {
			return "", err
		}
	} else {
		if err := s.local.SavePNG(h, data); err != nil {
			return "", err
		}
	}

	variants, err := s.resizer.Generate(data)
	if err != nil {
		return "", err
	}

	if variants.OG != nil {
		if useS3 {
			if err := s.s3.SaveOG(h, variants.OG); err != nil {
				return "", err
			}
		} else {
			if err := s.local.SaveOG(h, variants.OG); err != nil {
				return "", err
			}
		}
	}

	if variants.RE != nil {
		if useS3 {
			if err := s.s3.SaveRE(h, variants.RE); err != nil {
				return "", err
			}
		} else {
			if err := s.local.SaveRE(h, variants.RE); err != nil {
				return "", err
			}
		}
	}

	return h, nil
}
