package auth

import (
	"github.com/intezya/pkglib/crypto"
	"github.com/intezya/pkglib/logger"
)

const (
	defaultSaltLength = 32
)

type HashHelper struct{}

func NewHashHelper() *HashHelper {
	return &HashHelper{}
}

func (h *HashHelper) preHash(raw string) (prehash string) {
	return crypto.HashSHA256(raw)
}

func (h *HashHelper) EncodePassword(raw string) (hash string) {
	salt, err := crypto.Salt(defaultSaltLength)
	if err != nil {
		logger.Log.Errorf("Unexpected error while generating salt: %h", err)

		for {
			salt, err = crypto.Salt(defaultSaltLength)
			if err == nil {
				break
			}
		}
	}

	return crypto.HashArgon2(h.preHash(raw), salt, nil)
}

func (h *HashHelper) VerifyPassword(raw, hash string) bool {
	return crypto.VerifyArgon2(h.preHash(raw), hash, nil)
}

func (h *HashHelper) EncodeHardwareID(raw string) (hash string) {
	return crypto.HashSHA256(raw)
}

func (h *HashHelper) VerifyHardwareID(raw, hash string) bool {
	return h.EncodeHardwareID(raw) == hash
}

func (h *HashHelper) VerifyTokenHardwareID(tokenHash, userHash string) bool {
	return tokenHash == userHash
}
