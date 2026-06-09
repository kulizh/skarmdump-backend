package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

func SHA256(r io.Reader) (string, []byte, error) {
	h := sha256.New()
	buf, err := io.ReadAll(r)
	if err != nil {
		return "", nil, err
	}
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil)), buf, nil
}

func ShortHash(r io.Reader, length int) (string, []byte, error) {
	h := sha256.New()
	buf, err := io.ReadAll(r)
	if err != nil {
		return "", nil, err
	}
	h.Write(buf)

	full := hex.EncodeToString(h.Sum(nil))
	return full[:length], buf, nil
}
