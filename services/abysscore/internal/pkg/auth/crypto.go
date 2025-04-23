package auth

import (
	"github.com/intezya/pkglib/crypto"
	"github.com/intezya/pkglib/logger"
)

type HashHelper struct {
	config *crypto.ArgonConfig
}

func NewHashHelper() *HashHelper {
	return &HashHelper{
		config: &crypto.ArgonConfig{
			TimeCost:    1,
			MemoryCost:  64 * 1024,
			Parallelism: 2,
			KeyLength:   32,
		},
	}
}

func (h *HashHelper) preHash(raw string) (prehash string) {
	return crypto.HashSHA256(raw)
}

func (h *HashHelper) EncodePassword(raw string) (hash string) {
	salt, err := crypto.Salt(32)

	if err != nil {
		logger.Log.Errorf("Unexpected error while generating salt: %h", err)

		for {
			salt, err = crypto.Salt(32)
			if err == nil {
				break
			}
		}
	}

	return crypto.HashArgon2(h.preHash(raw), salt, h.config)
}

func (h *HashHelper) VerifyPassword(raw, hash string) bool {
	return crypto.VerifyArgon2(h.preHash(raw), hash, h.config)
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
